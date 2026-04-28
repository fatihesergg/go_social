package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseFollowStore interface {
	GetFollowerByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
	GetFollowingByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
	FollowUser(ctx context.Context, model model.Follow) error
	UnFollowUser(ctx context.Context, model model.Follow) error
}

type FollowStore struct {
	db *sql.DB
}

func NewFollowStore(db *sql.DB) BaseFollowStore {
	return &FollowStore{db: db}
}

func (s FollowStore) GetFollowerByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE follow_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting follower by userid: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		follows = append(follows, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}

func (s FollowStore) GetFollowingByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE user_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}
		return nil, fmt.Errorf("Error while getting followings by userid: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)

	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}

func (s FollowStore) FollowUser(ctx context.Context, model model.Follow) error {
	query := "INSERT INTO follows (id,user_id, follow_id) VALUES ($1, $2,$3)"
	_, err := s.db.ExecContext(ctx, query, model.ID, model.UserID, model.FollowID)
	if err != nil {
		return fmt.Errorf("Error while inserting follow: %w", err)
	}
	return nil
}
func (s FollowStore) UnFollowUser(ctx context.Context, model model.Follow) error {
	query := "DELETE FROM follows WHERE user_id = $1 AND follow_id = $2"
	result, err := s.db.ExecContext(ctx, query, model.UserID, model.FollowID)
	if err != nil {
		return fmt.Errorf("Error while deleting follow: %w", err)
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
