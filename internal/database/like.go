package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseLikeStore interface {
	LikePost(like *model.PostLike) error
	LikeComment(like *model.CommentLike) error
	UnlikePost(postID uuid.UUID, userID uuid.UUID) error
	UnlikeComment(commentID uuid.UUID, userID uuid.UUID) error
	IsPostLiked(postID uuid.UUID, userID uuid.UUID) (bool, error)
	IsCommentLiked(commentID uuid.UUID, userID uuid.UUID) (bool, error)
}

type LikeStore struct {
	DB *sql.DB
}

func NewLikeStore(db *sql.DB) BaseLikeStore {
	return &LikeStore{DB: db}
}

func (s *LikeStore) LikePost(like *model.PostLike) error {
	query := `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2)`
	_, err := s.DB.Exec(query, like.PostID, like.UserID)
	return err
}

func (s *LikeStore) LikeComment(like *model.CommentLike) error {
	query := `INSERT INTO comment_likes (comment_id, user_id) VALUES ($1, $2)`
	_, err := s.DB.Exec(query, like.CommentID, like.UserID)
	return err
}
func (s *LikeStore) UnlikePost(postID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
	_, err := s.DB.Exec(query, postID, userID)
	return err
}
func (s *LikeStore) UnlikeComment(commentID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`
	_, err := s.DB.Exec(query, commentID, userID)
	return err
}
func (s *LikeStore) IsPostLiked(postID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS ( SELECT  1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	err := s.DB.QueryRow(query, postID, userID).Scan(&result)
	if err != nil {
		return false, err
	}
	return result, nil
}
func (s *LikeStore) IsCommentLiked(commentID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`
	err := s.DB.QueryRow(query, commentID, userID).Scan(&result)
	if err != nil {
		return false, err
	}
	return result, nil
}
