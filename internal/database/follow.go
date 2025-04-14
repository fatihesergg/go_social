package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
)

type BaseFollowStore interface {
	GetFollowerByUserID(userID int64) ([]model.Follow, error)
	GetFollowingByUserID(userID int64) ([]model.Follow, error)
	FollowUser(userID, followID int64) error
	UnFollowUser(userID, followID int64) error
}

type FollowStore struct {
	db *sql.DB
}

func NewFollowStore(db *sql.DB) BaseFollowStore {
	return &FollowStore{db: db}
}

func (s FollowStore) GetFollowerByUserID(userID int64) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE user_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return follows, nil
}

func (s FollowStore) GetFollowingByUserID(userID int64) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE follow_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			return nil, err
		}
		follows = append(follows, follow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return follows, nil
}

func (s FollowStore) FollowUser(userID, followID int64) error {
	query := "INSERT INTO follows (user_id, follow_id) VALUES ($1, $2)"
	_, err := s.db.Exec(query, userID, followID)
	return err
}
func (s FollowStore) UnFollowUser(userID, followID int64) error {
	query := "DELETE FROM follows WHERE user_id = $1 AND follow_id = $2"
	_, err := s.db.Exec(query, userID, followID)
	return err
}
