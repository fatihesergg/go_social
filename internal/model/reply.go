package model

import "github.com/google/uuid"

type Reply struct {
	ID        uuid.UUID `json:"id"`
	CommentID uuid.UUID `json:"comment_id"`
	UserID    uuid.UUID `json:"-"`
	Message   string    `json:"message"`
	User      User      `json:"user"`
}
