package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/imrishuroy/legal-referral-notification-service/api"
	db "github.com/imrishuroy/legal-referral-notification-service/db/sqlc"
	"github.com/imrishuroy/legal-referral-notification-service/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Create a counter metric
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route"},
	)
)

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(requestsTotal)
}

func main() {

	log.Info().Msg("Welcome to LegalReferral Notification Service")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("cannot load config")
	}

	// db connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		fmt.Println("cannot connect to db:", err)
	}
	defer connPool.Close() // close db connection

	store := db.NewStore(connPool)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
	}

	go func() {
		err := api.ConnectConsumer(server)
		if err != nil {
			log.Error().Err(err).Msg("cannot connect consumer")
			panic(err)
		}
	}()

	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: nil,
	}

	go func() {
		http.HandleFunc("/health", healthCheck)
		http.HandleFunc("/", healthCheck)
		fmt.Println("Starting server at " + config.ServerAddress)
		log.Info().Msg("server address: " + config.ServerAddress)
		// Expose the /metrics endpoint for Prometheus to scrape
		http.Handle("/metrics", promhttp.Handler())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("cannot start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "OK")
}
