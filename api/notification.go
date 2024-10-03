package api

import (
	"context"
	"firebase.google.com/go/messaging"
	"github.com/rs/zerolog/log"
)

func (server *Server) sendNotification(userID string) error {

	// Get the device token from the database
	deviceToken, err := server.store.GetDeviceTokenByUserId(context.Background(), userID)

	if err != nil {
		log.Error().Err(err).Msg("error getting device token")
		return err
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "19:26",
		},
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implement push notification",
		},
		Token: deviceToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := server.client.Send(context.Background(), message)
	if err != nil {
		log.Error().Err(err).Msg("error sending message")
		return nil
	}
	// Response is a message ID string.
	log.Printf("Successfully sent message: %v\n", response)
	return nil
}
