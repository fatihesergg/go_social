package dto

import (
	"encoding/json"

	"github.com/fatihesergg/go_social/internal/model"
)

type NotificationResponse struct {
	Payload map[string]interface{}
	Message string
}

func NewNotificationResponse(model []model.Notification) []NotificationResponse {
	notifications := []NotificationResponse{}
	for _, notify := range model {
		payloadMap := map[string]interface{}{}
		_ = json.Unmarshal(notify.Payload, payloadMap)
		response := NotificationResponse{
			Payload: payloadMap,
			Message: notify.Message,
		}
		notifications = append(notifications, response)
	}
	return notifications
}
