package model

import (
	"database/sql"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID      `json:"id"`
	Content   string         `json:"content"`
	UserID    uuid.UUID      `json:"-"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	User      User           `json:"user"`
	Comments  []Comment      `json:"comments"`
}
