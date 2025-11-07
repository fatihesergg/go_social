package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseLikeStore interface {
	LikePost(like *model.PostLike) error
	LikeComment(like *model.CommentLike) error
	LikeReply(like *model.ReplyLike) error
	UnlikePost(postID uuid.UUID, userID uuid.UUID) error
	UnlikeComment(commentID uuid.UUID, userID uuid.UUID) error
	UnlikeReply(replyID uuid.UUID, userID uuid.UUID) error
	IsPostLiked(postID uuid.UUID, userID uuid.UUID) (bool, error)
	IsCommentLiked(commentID uuid.UUID, userID uuid.UUID) (bool, error)
	IsReplyLiked(replyID uuid.UUID, userID uuid.UUID) (bool, error)
}

type LikeStore struct {
	DB     *sql.DB
	logger *zap.Logger
}

func NewLikeStore(db *sql.DB, logger *zap.Logger) BaseLikeStore {
	return &LikeStore{DB: db, logger: logger.Named("like_store")}
}

func (s *LikeStore) LikePost(like *model.PostLike) error {
	query := `INSERT INTO post_likes (id,post_id, user_id) VALUES ($1, $2,$3)`
	_, err := s.DB.Exec(query, like.ID, like.PostID, like.UserID)
	if err != nil {
		s.logger.Error("Error while inserting post_like", zap.Error(err))
		return err
	}
	return nil
}

func (s *LikeStore) LikeComment(like *model.CommentLike) error {
	query := `INSERT INTO comment_likes (id,comment_id, user_id) VALUES ($1, $2,$3)`
	_, err := s.DB.Exec(query, like.ID, like.CommentID, like.UserID)
	if err != nil {
		s.logger.Error("Error while inserting comment_like", zap.Error(err))
		return err
	}
	return nil
}
func (s *LikeStore) UnlikePost(postID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`
	result, err := s.DB.Exec(query, postID, userID)
	if err != nil {
		s.logger.Error("Error while delete post_like", zap.Error(err))
		return err
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

func (s *LikeStore) LikeReply(like *model.ReplyLike) error {
	query := `INSERT INTO reply_likes  (id,reply_id,user_id) VALUES ($1,$2,$3)`
	_, err := s.DB.Exec(query, like.ID, like.ReplyID, like.UserID)
	if err != nil {
		s.logger.Error("Error while inserting reply_like", zap.Error(err))
		return err
	}
	return nil
}

func (s *LikeStore) UnlikeComment(commentID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`
	result, err := s.DB.Exec(query, commentID, userID)
	if err != nil {
		s.logger.Error("Error while deleting comment_like", zap.Error(err))
		return err
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
func (s *LikeStore) UnlikeReply(replyID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM reply_likes WHERE reply_id = $1 AND user_id = $2`
	result, err := s.DB.Exec(query, replyID, userID)
	if err != nil {
		s.logger.Error("Error while deleting reply_like", zap.Error(err))
		return err
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

func (s *LikeStore) IsPostLiked(postID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS ( SELECT  1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`
	err := s.DB.QueryRow(query, postID, userID).Scan(&result)
	if err != nil {
		s.logger.Error("Error while checking post liked", zap.Error(err))
		return false, err
	}
	return result, nil
}
func (s *LikeStore) IsCommentLiked(commentID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`
	err := s.DB.QueryRow(query, commentID, userID).Scan(&result)
	if err != nil {
		s.logger.Error("Error while checking comment liked", zap.Error(err))
		return false, err
	}
	return result, nil
}
func (s *LikeStore) IsReplyLiked(replyID uuid.UUID, userID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS (SELECT 1 FROM reply_likes WHERE reply_id = $1 AND user_id = $2)`
	err := s.DB.QueryRow(query, replyID, userID).Scan(&result)
	if err != nil {
		s.logger.Error("Error while checking reply liked", zap.Error(err))
		return false, err
	}
	return result, nil
}
