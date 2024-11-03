package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/imrishuroy/legal-referral-notification-service/util"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

type Notification struct {
	UserID           string `json:"user_id" binding:"required"`
	SenderID         string `json:"sender_id" binding:"required"`
	TargetID         string `json:"target_id" binding:"required"`
	TargetType       string `json:"target_type" binding:"required"`
	NotificationType string `json:"notification_type" binding:"required"`
	AlreadyLiked     string `json:"already_liked"`
}

func ConnectConsumer(server *Server) error {

	topic := "likes"
	worker, err := createConsumer(server.config)
	if err != nil {
		return err
	}
	// Calling ConsumePartition. It will open one connection per broker
	// and share it for all partitions that live on it.
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest) // sarama.OffsetOldest
	if err != nil {
		return err
	}
	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed
	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error().Err(err).Msg("error")
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, msg.Topic, string(msg.Key))

				var notification Notification

				// Unmarshal the JSON into the struct
				err := json.Unmarshal(msg.Value, &notification)
				if err != nil {
					log.Error().Err(err).Msg("cannot unmarshal data")
				}
				err = server.processNotification(notification)
				if err != nil {
					log.Error().Err(err).Msg("cannot post to news feed")
				}
			case <-sigchan:
				log.Info().Msg("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Info().Msgf("Processed %d messages", msgCount)

	if err := worker.Close(); err != nil {
		return err
	}

	return nil
}

func createConsumer(config util.Config) (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext //"PLAIN"
	cfg.Net.SASL.Enable = true

	cfg.Net.SASL.User = config.SASLUsername
	cfg.Net.SASL.Password = config.SASLPassword

	tlsConfig := tls.Config{}
	cfg.Net.TLS.Enable = true
	cfg.Net.TLS.Config = &tlsConfig

	// Create new consumer
	conn, err := sarama.NewConsumer([]string{config.BootStrapServers}, cfg)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
