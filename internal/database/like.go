package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseLikeStore interface {
	LikePost(ctx context.Context, like *model.PostLike) error
	LikeComment(ctx context.Context, like *model.CommentLike) error
	LikeReply(ctx context.Context, like *model.ReplyLike) error
	UnlikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	UnlikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
	UnlikeReply(ctx context.Context, replyID uuid.UUID, userID uuid.UUID) error
	IsPostLiked(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (bool, error)
	IsCommentLiked(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (bool, error)
	IsReplyLiked(ctx context.Context, replyID uuid.UUID, userID uuid.UUID) (bool, error)
}

type LikeStore struct {
	DB *sql.DB
}

func NewLikeStore(db *sql.DB) BaseLikeStore {
	return &LikeStore{DB: db}
}

func (s *LikeStore) LikePost(ctx context.Context, like *model.PostLike) error {
	query := `INSERT INTO post_likes (id,post_id, user_id) VALUES ($1, $2,$3)`
	_, err := s.DB.ExecContext(ctx, query, like.ID, like.PostID, like.UserID)
	if err != nil {
		return fmt.Errorf("error while inserting post_like: %w", err)
	}
	return nil
}

func (s *LikeStore) LikeComment(ctx context.Context, like *model.CommentLike) error {
	query := `INSERT INTO comment_likes (id,comment_id, user_id) VALUES ($1, $2,$3)`
	_, err := s.DB.ExecContext(ctx, query, like.ID, like.CommentID, like.UserID)
	if err != nil {
		return fmt.Errorf("error while inserting comment_like: %w", err)
	}
	return nil
}
func (s *LikeStore) UnlikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
	result, err := s.DB.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("error while delete post_like: %w", err)
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

func (s *LikeStore) LikeReply(ctx context.Context, like *model.ReplyLike) error {
	query := `INSERT INTO reply_likes  (id,reply_id,user_id) VALUES ($1,$2,$3)`
	_, err := s.DB.ExecContext(ctx, query, like.ID, like.ReplyID, like.UserID)
	if err != nil {
		return fmt.Errorf("error while inserting reply_like: %w", err)
	}
	return nil
}

func (s *LikeStore) UnlikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`
	result, err := s.DB.ExecContext(ctx, query, commentID, userID)
	if err != nil {
		return fmt.Errorf("error while deleting comment_like: %w", err)
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
func (s *LikeStore) UnlikeReply(ctx context.Context, replyID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM reply_likes WHERE reply_id = $1 AND user_id = $2`
	result, err := s.DB.ExecContext(ctx, query, replyID, userID)
	if err != nil {
		return fmt.Errorf("error while deleting reply_like: %w", err)
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

func (s *LikeStore) IsPostLiked(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS ( SELECT  1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	err := s.DB.QueryRowContext(ctx, query, postID, userID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("error while checking post liked: %w", err)
	}
	return result, nil
}
func (s *LikeStore) IsCommentLiked(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`
	err := s.DB.QueryRowContext(ctx, query, commentID, userID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("error while checking comment liked: %w", err)
	}
	return result, nil
}
func (s *LikeStore) IsReplyLiked(ctx context.Context, replyID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM reply_likes WHERE reply_id = $1 AND user_id = $2)`
	err := s.DB.QueryRowContext(ctx, query, replyID, userID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("error while checking reply liked: %w", err)
	}
	return result, nil
}
