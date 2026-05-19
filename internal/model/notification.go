package model

import (
	"time"

	"github.com/google/uuid"
)

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
	Timestamp         time.Time
}

type PostLikePayload struct {
	PostID     string
	PostUserID string
	LikerID    string
}

type PostLikeNotification struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Notification_type NotificationType
	Payload           any
	Message           string
	IsRead            bool
	Timestamp         time.Time
}
