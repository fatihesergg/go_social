package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
)

type NotificationResponse struct {
	Payload any
	Message string
}

func NewNotificationResponse(model []model.Notification) []NotificationResponse {
	notifications := []NotificationResponse{}
	for _, notify := range model {
		response := NotificationResponse{
			Payload: notify.Payload,
			Message: notify.Message,
		}
		notifications = append(notifications, response)
	}
	return notifications
}
