package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID      `json:"id"`
	PostID    int64          `json:"post_id"`
	UserID    int64          `json:"-"`
	Content   string         `json:"content"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	User      User           `json:"user"`
}
