package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseUserStore interface {
	GetUsersByUsername(ctx context.Context, userName string) ([]model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) BaseUserStore {
	return &UserStore{DB: db}
}

func (s *UserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id,name,last_name,username,email,password,avatar,created_at,updated_at FROM users WHERE id = $1"
	row := s.DB.QueryRowContext(ctx, query, id.String())

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("error while getting user by id: %w", err)
	}
	return user, nil
}

func (s *UserStore) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id,name,last_name,username,email,password,avatar,created_at,updated_at FROM users WHERE username = $1"
	row := s.DB.QueryRowContext(ctx, query, username)

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("error while getting user by username: %w", err)
	}
	return user, nil
}

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := "SELECT id, name, last_name, username, email, password, created_at, updated_at, avatar FROM users WHERE email = $1"
	row := s.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.ID, &user.Name, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Avatar)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("error while getting user by email: %w", err)
	}
	return user, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user *model.User) error {
	var id uuid.UUID
	query := "INSERT INTO users (id, name, last_name, username, email, password, avatar) VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id"
	err := s.DB.QueryRowContext(ctx, query, user.ID, user.Name, user.LastName, user.Username, user.Email, user.Password, user.Avatar).Scan(&id)
	if err != nil {
		return fmt.Errorf("error while creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET name = $1, last_name = $2, username = $3, email = $4, password = $5 WHERE id = $6"
	result, err := s.DB.ExecContext(ctx, query, user.Name, user.LastName, user.Username, user.Email, user.Password, user.ID.String())
	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *UserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"
	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error while deleting user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *UserStore) GetUsersByUsername(ctx context.Context, userName string) ([]model.User, error) {
	users := []model.User{}
	query := "SELECT id,name,last_name,username FROM users WHERE username ILIKE '%' || $1 || '%'"
	rows, err := s.DB.QueryContext(ctx, query, userName)

	if err != nil {
		return nil, fmt.Errorf("error while getting users by username: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		user := model.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("error while scanning result set: %w", err)
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error in result row: %w", err)
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return users, err

}
