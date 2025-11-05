package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseFollowStore interface {
	GetFollowerByUserID(userID uuid.UUID) ([]model.Follow, error)
	GetFollowingByUserID(userID uuid.UUID) ([]model.Follow, error)
	FollowUser(model model.Follow) error
	UnFollowUser(model model.Follow) error
}

type FollowStore struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewFollowStore(db *sql.DB, logger *zap.Logger) BaseFollowStore {
	return &FollowStore{db: db, logger: logger.Named("follow_store")}
}

func (s FollowStore) GetFollowerByUserID(userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE follow_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		s.logger.Error("Error while getting follower by userid", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			s.logger.Error("Error while scanning result", zap.Error(err))
			return nil, err
		}
		follows = append(follows, follow)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("Error in result row", zap.Error(err))
		return nil, err
	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}

func (s FollowStore) GetFollowingByUserID(userID uuid.UUID) ([]model.Follow, error) {
	var follows []model.Follow
	query := "SELECT id, user_id, follow_id FROM follows WHERE user_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return follows, nil
		}
		s.logger.Error("Error while getting followings by userid", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var follow model.Follow
		if err := rows.Scan(&follow.ID, &follow.UserID, &follow.FollowID); err != nil {
			s.logger.Error("Error while scanning result", zap.Error(err))
			return nil, err
		}
		follows = append(follows, follow)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("Error in result row", zap.Error(err))
		return nil, err

	}

	if len(follows) == 0 {
		return nil, sql.ErrNoRows
	}

	return follows, nil
}

func (s FollowStore) FollowUser(model model.Follow) error {
	query := "INSERT INTO follows (id,user_id, follow_id) VALUES ($1, $2,$3)"
	_, err := s.db.Exec(query, model.ID, model.UserID, model.FollowID)
	if err != nil {
		s.logger.Error("Error while inserting follow", zap.Error(err))
		return err
	}
	return nil
}
func (s FollowStore) UnFollowUser(model model.Follow) error {
	query := "DELETE FROM follows WHERE user_id = $1 AND follow_id = $2"
	_, err := s.db.Exec(query, model.UserID, model.FollowID)
	if err != nil {
		s.logger.Error("Error while deleting follow", zap.Error(err))
		return err
	}
	return nil
}
