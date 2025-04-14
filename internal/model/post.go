package model

import "database/sql"

type Post struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	UserID    int64          `json:"user_id"`
	Image     sql.NullString `json:"image"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	User      User           `json:"user"`
	Comments  []Comment      `json:"comments"`
}
