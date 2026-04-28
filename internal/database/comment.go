package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentStore interface {
	HasAccessToComment(ctx context.Context, userID, commentID uuid.UUID) (bool, error)
	GetCommentsByPostID(ctx context.Context, postID, userID uuid.UUID) ([]model.Comment, error)
	GetCommentByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	CreateComment(ctx context.Context, comment *model.Comment) error
	UpdateComment(ctx context.Context, comment *model.Comment) error
	DeleteComment(ctx context.Context, id uuid.UUID) error
}

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) BaseCommentStore {
	return &CommentStore{
		db: db,
	}
}

func (cs *CommentStore) HasAccessToComment(ctx context.Context, userID, commentID uuid.UUID) (bool, error) {
	var result bool
	query := `
	SELECT EXISTS(SELECT 1 FROM comments
	JOIN posts AS comments_post ON comments_post.id = comments.post_id
	LEFT JOIN follows ON follows.follow_id = comments_post.user_id
	WHERE comments.id = $2 
	AND (
	comments_post.visibility = 'public' 
	OR comments.user_id = $1
	OR (follows.user_id IS NOT NULL AND follows.user_id = $1)
	))`

	err := cs.db.QueryRowContext(ctx, query, userID, commentID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("Error while checking comment access: %w", err)
	}
	return result, nil

}

func (cs CommentStore) GetCommentsByPostID(ctx context.Context, postID, userID uuid.UUID) ([]model.Comment, error) {
	var comments []model.Comment
	query := `
	WITH comment_likes_count AS (
	SELECT comment_id,COUNT(*) AS likes_count FROM comment_likes
	GROUP BY comment_id
	),

	reply_count AS (
	SELECT comment_id,COUNT(*) AS replies_count FROM replies
	GROUP BY comment_id
	),

	user_likes AS (
	SELECT comment_id FROM comment_likes
	WHERE user_id = $2 
	),

	user_follows AS (
	SELECT follow_id FROM follows
	WHERE user_id = $2
	)

	
	SELECT
	comments.id,
	comments.post_id,
	comments.content,
	comments.created_at,
	comments.updated_at,

	users.id,
	users.name,
	users.last_name,
	users.username,


	COALESCE(comment_likes_count.likes_count ,0) AS total_likes,
	COALESCE(reply_count.replies_count,0) AS total_reply,

	(user_likes.comment_id IS NOT NULL) AS is_liked,
	(user_follows.follow_id IS NOT NULL) AS is_following


	FROM comments JOIN users ON comments.user_id = users.id
	LEFT JOIN comment_likes_count ON comment_likes_count.comment_id = comments.id
	LEFT JOIN reply_count ON reply_count.comment_id = comments.id
	LEFT JOIN user_likes ON user_likes.comment_id = comments.id
	LEFT JOIN user_follows ON user_follows.follow_id = users.id
	WHERE comments.post_id = $1`
	rows, err := cs.db.QueryContext(ctx, query, postID, userID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting comments by postid: %w", err)
	}
	defer rows.Close()
	commentMap := make(map[uuid.UUID]*model.Comment)
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.ID, &comment.User.Name, &comment.User.LastName, &comment.User.Username,
			&comment.LikeCount, &comment.ReplyCount,
			&comment.IsLiked, &comment.IsFollowing)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		if _, ok := commentMap[comment.ID]; !ok {
			commentMap[comment.ID] = &comment
		}

	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}

	for _, comment := range commentMap {
		comments = append(comments, *comment)
	}

	if len(comments) == 0 {
		return nil, sql.ErrNoRows
	}

	return comments, nil

}
func (cs CommentStore) GetCommentByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	query := `SELECT comments.id,comments.post_id,comments.user_id,comments.content,comments.created_at,comments.updated_at,
		users.name,users.last_name,users.username,users.email	
	FROM comments 
	JOIN users ON comments.user_id = users.id
	WHERE comments.id = $1`
	err := cs.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email)

	if err != nil {
		return nil, fmt.Errorf("Error wile getting comment by id: %w", err)
	}

	return &comment, nil
}

func (cs CommentStore) CreateComment(ctx context.Context, comment *model.Comment) error {

	query := "INSERT INTO comments (id,post_id, user_id, content) VALUES ($1, $2, $3,$4)"
	_, err := cs.db.ExecContext(ctx, query, comment.ID, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return fmt.Errorf("Error while inserting comment: %w", err)
	}
	return nil
}

func (cs CommentStore) UpdateComment(ctx context.Context, comment *model.Comment) error {
	query := "UPDATE comments SET content = $1 WHERE id = $2"
	result, err := cs.db.ExecContext(ctx, query, comment.Content, comment.ID)
	if err != nil {
		return fmt.Errorf("Error while updating comment: %w", err)
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

func (cs CommentStore) DeleteComment(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM comments WHERE id = $1"
	result, err := cs.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Error while deleting comment: %w", err)
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
