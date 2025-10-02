package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseUserStore interface {
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user model.User) error
	UpdateUser(user model.User) error
	DeleteUser(id uuid.UUID) error
}

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) BaseUserStore {
	return &UserStore{DB: db}
}

func (s *UserStore) GetUserByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id,name,last_name,username,email,password,avatar,created_at,updated_at FROM users WHERE id = $1"
	row := s.DB.QueryRow(query, id.String())

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStore) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id,name,last_name,username,email,password,avatar,created_at,updated_at FROM users WHERE username = $1"
	row := s.DB.QueryRow(query, username)

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStore) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id, name, last_name, username, email, password, created_at, updated_at, avatar FROM users WHERE email = $1"
	row := s.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Avatar)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *UserStore) CreateUser(user model.User) error {
	var id uuid.UUID
	query := "INSERT INTO users (name, last_name, username, email, password, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err := s.DB.QueryRow(query, user.Name, user.LastName, user.Username, user.Email, user.Password, user.Avatar).Scan(&id)
	if err != nil {

		return err
	}
	return nil
}

func (s *UserStore) UpdateUser(user model.User) error {
	query := "UPDATE users SET name = $1, last_name = $2, username = $3, email = $4, password = $5, avatar = $6 WHERE id = $7"
	_, err := s.DB.Exec(query, user.Name, user.LastName, user.Username, user.Email, user.Password, user.Avatar, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
