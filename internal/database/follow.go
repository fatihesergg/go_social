package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseFollowStore interface {
	GetFollowerByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
	GetFollowingByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error)
	UpsertFollowRequest(ctx context.Context, followModel model.Follow) (string, error)
	DeleteFollowRequest(ctx context.Context, model model.Follow, status model.FollowReqStatus) (string, string, error)
	IsFollowing(ctx context.Context, model model.Follow) (bool, error)
	UpdateFollowStatus(ctx context.Context, model model.Follow) (string, string, error)
}

type FollowStore struct {
	db *sql.DB
}

func NewFollowStore(db *sql.DB) BaseFollowStore {
	return &FollowStore{db: db}
}

func (s *FollowStore) GetFollowerByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE follow_id = $1 AND status = 'accepted' "
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting follower by userid: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, fmt.Errorf("error while scanning result: %w", err)
		}
		follows = append(follows, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in result row: %w", err)
	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}

func (s *FollowStore) GetFollowingByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE user_id = $1 AND status = 'accepted' "
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return follows, nil
		}
		return nil, fmt.Errorf("error while getting followings by userid: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, fmt.Errorf("error while scanning result: %w", err)
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in result row: %w", err)

	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}
func (s *FollowStore) UpsertFollowRequest(ctx context.Context, followModel model.Follow) (string, error) {
	query := `
	INSERT INTO follows(id,user_id,follow_id,status) VALUES($1,$2,$3,$4)
	ON CONFLICT (user_id,follow_id) 
	DO UPDATE SET 
	status = CASE WHEN follows.status = 'rejected' THEN 'pending' ELSE follows.status END,
	updated_at = CASE WHEN follows.status = 'rejected' THEN NOW() ELSE follows.updated_at END
	RETURNING old.status
	`
	var old *string

	row := s.db.QueryRowContext(ctx, query, followModel.ID, followModel.UserID, followModel.FollowID, followModel.Status)

	if err := row.Scan(&old); err != nil {
		return "", fmt.Errorf("error while scanning upsert follow query: %w", err)
	}

	if old == nil {
		return "", nil
	}

	return *old, nil

}

func (s *FollowStore) UpdateFollowStatus(ctx context.Context, model model.Follow) (string, string, error) {
	query := `WITH old_status AS (
	SELECT status FROM follows WHERE user_id = $2 AND follow_id = $3
	),

	new_status AS (
		UPDATE follows SET status = $1 WHERE user_id = $2 AND follow_id = $3 AND status = 'pending'
		RETURNING status
	)

	SELECT (SELECT status FROM old_status),(SELECT status FROM new_status)
	
	`
	row := s.db.QueryRowContext(ctx, query, model.Status, model.UserID, model.FollowID)

	var old *string
	var new *string
	if err := row.Scan(&old, &new); err != nil {
		return "", "", fmt.Errorf("error while scanning result: %w", err)
	}
	if old == nil && new == nil {
		return "", "", nil
	}
	if old == nil && new != nil {
		return "", *new, nil
	}
	if old != nil && new == nil {
		return *old, "", nil
	}

	return *old, *new, nil
}

func (s *FollowStore) DeleteFollowRequest(ctx context.Context, model model.Follow, status model.FollowReqStatus) (string, string, error) {
	query := `WITH old_status AS (
	SELECT status FROM follows WHERE user_id = $2 AND follow_id = $3
	),

	new_status AS (
		DELETE FROM follows WHERE user_id = $2 AND follow_id = $3 AND status = $1
		RETURNING status
	)

	SELECT (SELECT status FROM old_status),(SELECT status FROM new_status)`

	var old *string
	var new *string
	row := s.db.QueryRowContext(ctx, query, status, model.UserID, model.FollowID)

	if err := row.Scan(&old, &new); err != nil {
		return "", "", fmt.Errorf("error while scanning result: %w", err)
	}
	if old == nil && new == nil {
		return "", "", nil
	}
	if old == nil && new != nil {
		return "", *new, nil
	}
	if old != nil && new == nil {
		return *old, "", nil
	}

	return *old, *new, nil
}

func (s *FollowStore) IsFollowing(ctx context.Context, model model.Follow) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM follows WHERE user_id = $1 AND follow_id = $2 AND status = 'accepted' )"
	result := s.db.QueryRowContext(ctx, query, model.UserID, model.FollowID)

	var check bool
	if err := result.Scan(&check); err != nil {
		return false, err
	}
	return check, nil
}
