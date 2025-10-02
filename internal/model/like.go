package model

import "github.com/google/uuid"

type PostLike struct {
	ID     uuid.UUID `json:"id"`
	PostID uuid.UUID `json:"post_id"`
	UserID uuid.UUID `json:"user_id"`
}
type CommentLike struct {
	ID        uuid.UUID `json:"id"`
	CommentID uuid.UUID `json:"comment_id"`
	UserID    uuid.UUID `json:"user_id"`
}
