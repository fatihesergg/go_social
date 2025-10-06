package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentStore interface {
	GetCommentsByPostID(postID uuid.UUID) ([]model.Comment, error)
	GetCommentByID(id uuid.UUID) (*model.Comment, error)
	CreateComment(comment *model.Comment) error
	UpdateComment(comment *model.Comment) error
	DeleteComment(id uuid.UUID) error
}

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) BaseCommentStore {
	return &CommentStore{
		db: db,
	}
}

func (cs CommentStore) GetCommentsByPostID(postID uuid.UUID) ([]model.Comment, error) {
	var comments []model.Comment
	query := `SELECT comments.id,comments.post_id,comments.user_id,comments.content,comments.image,comments.created_at,comments.updated_at,
	users.name,users.last_name,users.username,users.email
	FROM comments JOIN users ON comments.user_id = users.id WHERE comments.post_id = $1`
	rows, err := cs.db.Query(query, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil

}
func (cs CommentStore) GetCommentByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	query := `SELECT comments.id,comments.post_id,comments.user_id,comments.content,comments.image,comments.created_at,comments.updated_at,
		users.name,users.last_name,users.username,users.email	
	FROM comments 
	JOIN users ON comments.user_id = users.id
	WHERE comments.id = $1`
	err := cs.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (cs CommentStore) CreateComment(comment *model.Comment) error {

	query := "INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)"
	_, err := cs.db.Exec(query, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return err
	}
	return nil
}

func (cs CommentStore) UpdateComment(comment *model.Comment) error {
	query := "UPDATE comments SET content = $1 WHERE id = $2"
	_, err := cs.db.Exec(query, comment.Content, comment.ID)
	if err != nil {
		return err
	}
	return nil
}

func (cs CommentStore) DeleteComment(id uuid.UUID) error {
	query := "DELETE FROM comments WHERE id = $1"
	_, err := cs.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
