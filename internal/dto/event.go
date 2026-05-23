package dto

import "github.com/google/uuid"

type PostLikedEvent struct {
	LikerID uuid.UUID `json:"liker_id"`
	PostID  uuid.UUID `json:"post_id"`
}
type CommentLikedEvent struct {
	LikerID   uuid.UUID `json:"liker_id"`
	CommentID uuid.UUID `json:"post_id"`
}

type ReplyLikedEvent struct {
	LikerID uuid.UUID `json:"liker_id"`
	ReplyID uuid.UUID `json:"reply_id"`
}
