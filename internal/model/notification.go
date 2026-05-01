package model

import "github.com/google/uuid"

type NotificationType int

const (
	PostLikedNotification NotificationType = iota
)

var NotificationTypeStr = map[NotificationType]string{
	PostLikedNotification: "post_liked",
}

func (nt NotificationType) String() string {
	return NotificationTypeStr[nt]
}

type Notification struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Notification_type NotificationType
	Payload           []byte
	Message           string
	IsRead            bool
}
