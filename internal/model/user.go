package model

import "time"

type User struct {
	ID        int64     `json:"-"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Avatar    *string   `json:"avatar"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
