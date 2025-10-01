package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
)

type BaseLikeStore interface {
	CreateLike(like model.Like) error
	DeleteLike(postID int64, userID int64) error
	IsLiked(postID int64, userID int64) (bool, error)
}

type LikeStore struct {
	DB *sql.DB
}

func NewLikeStore(db *sql.DB) BaseLikeStore {
	return &LikeStore{DB: db}
}

func (s *LikeStore) CreateLike(like model.Like) error {
	query := `INSERT INTO likes (post_id, user_id) VALUES ($1, $2)`
	_, err := s.DB.Exec(query, like.PostID, like.UserID)
	return err
}

func (s *LikeStore) DeleteLike(postID int64, userID int64) error {
	query := `DELETE FROM likes WHERE post_id = $1 AND user_id = $2`
	_, err := s.DB.Exec(query, postID, userID)
	return err
}

func (s *LikeStore) IsLiked(postID int64, userID int64) (bool, error) {
	var like model.Like
	query := `SELECT * FROM likes WHERE post_id = $1 AND user_id = $2`
	err := s.DB.QueryRow(query, postID, userID).Scan(&like.ID, &like.PostID, &like.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
