package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType int

const (
	PostLikeNotificationType    NotificationType = iota
	CommentLikeNotificationType NotificationType = iota
	ReplyLikeNotificationType   NotificationType = iota
)

var NotificationTypeStr = map[NotificationType]string{
	PostLikeNotificationType:    "post_like",
	CommentLikeNotificationType: "comment_like",
	ReplyLikeNotificationType:   "reply_like",
}

func (nt NotificationType) String() string {
	return NotificationTypeStr[nt]
}

type Notification struct {
	ID                uuid.UUID        `json:"id"`
	UserID            uuid.UUID        `json:"user_id"`
	Notification_type NotificationType `json:"notification_type"`
	Payload           []byte           `json:"payload"`
	Message           string           `json:"message"`
	IsRead            bool             `json:"is_read"`
	Timestamp         time.Time        `json:"timestamp"`
}

type PostLikePayload struct {
	PostID     string `json:"post_id"`
	PostUserID string `json:"post_user_id"`
	LikerID    string `json:"liker_id"`
}

type PostLikeNotification struct {
	ID                uuid.UUID        `json:"id"`
	UserID            uuid.UUID        `json:"user_id"`
	Notification_type NotificationType `json:"notification_type"`
	Payload           any              `json:"payload"`
	Message           string           `json:"message"`
	IsRead            bool             `json:"is_read"`
	Timestamp         time.Time        `json:"timestamp"`
}

type CommentLikeNotification struct {
	ID                uuid.UUID        `json:"id"`
	UserID            uuid.UUID        `json:"user_id"`
	Notification_type NotificationType `json:"notification_type"`
	Payload           any              `json:"payload"`
	Message           string           `json:"message"`
	IsRead            bool             `json:"is_read"`
	Timestamp         time.Time        `json:"timestamp"`
}

type CommentLikePayload struct {
	LikerID       uuid.UUID `json:"liker_id"`
	PostID        uuid.UUID `json:"post_id"`
	CommentID     uuid.UUID `json:"comment_id"`
	CommentUserID uuid.UUID `json:"comment_user_id"`
}

type ReplyLikeNotification struct {
	ID               uuid.UUID        `json:"id"`
	UserID           uuid.UUID        `json:"user_id"`
	NotificationType NotificationType `json:"notification_type"`
	Message          string           `json:"message"`
	Payload          ReplyLikePayload `json:"payload"`
	IsRead           bool             `json:"is_read"`
	Timestamp        time.Time        `json:"timestamp"`
}
type ReplyLikePayload struct {
	ReplyID uuid.UUID `json:"reply_id"`
	LikerID uuid.UUID `json:"liker_id"`
}
