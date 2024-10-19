package api

import (
	"context"
	"firebase.google.com/go/messaging"
	"fmt"
	db "github.com/imrishuroy/legal-referral-notification-service/db/sqlc"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (server *Server) processNotification(notification Notification) error {

	// getting author details
	res, err := server.store.GetUserNameByUserId(context.Background(), notification.SenderID)

	if err != nil {
		log.Error().Err(err).Msg("error getting user name")
		return err
	}

	name := res.FirstName + " " + res.LastName

	notificationMsg := notificationMsg(name, notification.TargetType, notification.NotificationType)

	// create notification**

	targetID, err := strconv.Atoi(notification.TargetID)

	notificationReq := createNotificationReq{
		UserID:           notification.UserID,
		SenderID:         notification.SenderID,
		TargetID:         int32(targetID),
		TargetType:       notification.TargetType,
		NotificationType: notification.NotificationType,
		Message:          notificationMsg,
	}

	err = server.createNotification(notificationReq)
	if err != nil {
		log.Error().Err(err).Msg("error creating notification")
		return err
	}

	// send notification**

	// Get the recipient's device token
	deviceToken, err := server.store.GetDeviceTokenByUserId(context.Background(), notification.UserID)

	if err != nil {
		log.Error().Err(err).Msg("error getting device token")
		return err
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"click_action":      "FLUTTER_NOTIFICATION_CLICK",
			"sender_id":         notification.SenderID,
			"target_id":         fmt.Sprintf("%d", notification.TargetID),
			"target_type":       notification.TargetType,
			"notification_type": notification.NotificationType,
		},
		Notification: &messaging.Notification{
			Title: "Legal Referral",
			Body:  notificationMsg,
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

type createNotificationReq struct {
	UserID           string
	SenderID         string
	TargetID         int32
	TargetType       string
	NotificationType string
	Message          string
}

func (server *Server) createNotification(req createNotificationReq) error {

	args := db.CreateNotificationParams{
		UserID:           req.UserID,
		SenderID:         req.SenderID,
		TargetID:         req.TargetID,
		TargetType:       req.TargetType,
		NotificationType: req.NotificationType,
		Message:          req.Message,
	}

	_, err := server.store.CreateNotification(context.Background(), args)
	if err != nil {
		log.Error().Err(err).Msg("error creating notification")
		return err
	}

	return nil
}

func notificationMsg(authorName, targetType, notificationType string) string {
	var message string

	// Create the notification message based on the type
	switch notificationType {
	case "like":
		message = fmt.Sprintf("%s liked your %s.", authorName, targetType)
	case "comment":
		message = fmt.Sprintf("%s commented on your %s.", authorName, targetType)
	case "share":
		message = fmt.Sprintf("%s shared your %s.", authorName, targetType)
	case "follow":
		message = fmt.Sprintf("%s started following you.", authorName)
	case "mention":
		message = fmt.Sprintf("%s mentioned you in a %s.", authorName, targetType)
	default:
		message = fmt.Sprintf("%s interacted with your %s.", authorName, targetType)
	}

	return message
}
