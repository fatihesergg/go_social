package model

import "github.com/google/uuid"

type Reply struct {
	ID          uuid.UUID `json:"id"`
	CommentID   uuid.UUID `json:"comment_id"`
	ParentID    uuid.UUID `json:"parent_id"`
	UserID      uuid.UUID `json:"-"`
	Message     string    `json:"message"`
	User        User      `json:"user"`
	ReplyCount  int       `json:"total_replies"`
	LikeCount   int       `json:"total_likes"`
	IsLiked     bool      `json:"is_liked"`
	IsFollowing bool      `json:"is_followed"`
}
