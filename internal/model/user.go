package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"-"`
	Password  string    `json:"-"`
	Avatar    *string   `json:"avatar"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
