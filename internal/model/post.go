package model

import (
	"github.com/google/uuid"
)

type Post struct {
	ID           uuid.UUID `json:"id"`
	Content      string    `json:"content"`
	UserID       uuid.UUID `json:"-"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	User         User      `json:"user"`
	LikeCount    int       `json:"total_likes"`
	CommentCount int       `json:"total_comment"`
	IsLiked      bool      `json:"is_liked"`
	IsFollowing  bool      `json:"is_followed"`
	Comments     []Comment `json:"comments"`
}
