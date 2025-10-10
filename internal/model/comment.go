package model

import (
	"github.com/google/uuid"
)

type Comment struct {
	ID          uuid.UUID `json:"id"`
	PostID      uuid.UUID `json:"-"`
	UserID      uuid.UUID `json:"-"`
	Content     string    `json:"content"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	User        User      `json:"user"`
	Replies     []Reply   `json:"replies"`
	LikeCount   int       `json:"total_likes"`
	ReplyCount  int       `json:"total_reply"`
	IsLiked     bool      `json:"is_liked"`
	IsFollowing bool      `json:"is_followed"`
}
