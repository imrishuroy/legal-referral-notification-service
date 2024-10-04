package api

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	db "github.com/imrishuroy/legal-referral-notification-service/db/sqlc"
	"github.com/imrishuroy/legal-referral-notification-service/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type Server struct {
	config util.Config
	store  db.Store
	client *messaging.Client
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	ctx := context.Background()
	opt := option.WithCredentialsFile("./service-account-key.json")
	log.Info().Msg("creating firebase app")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Error().Err(err).Msg("cannot create firebase app")
	}
	// Obtain a messaging.Client from the App.

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error getting Messaging client")
	}

	server := &Server{
		config: config,
		store:  store,
		client: client,
	}

	return server, nil
}
