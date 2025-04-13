package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentStore interface {
	GetCommentsByPostID(postID int64) ([]model.Comment, error)
	GetCommentByID(id uuid.UUID) (*model.Comment, error)
	CreateComment(comment model.Comment) error
	UpdateComment(comment model.Comment) error
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

func (cs CommentStore) GetCommentsByPostID(postID int64) ([]model.Comment, error) {
	var commets []model.Comment
	query := "SELECT id,post_id,user_id,content,image,created_at,updated_at FROM comments WHERE post_id = $1"
	rows, err := cs.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var commet model.Comment
		err := rows.Scan(&commet.ID, &commet.PostID, &commet.UserID, &commet.Content, &commet.Image, &commet.CreatedAt, &commet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		commets = append(commets, commet)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return commets, nil

}
func (cs CommentStore) GetCommentByID(id uuid.UUID) (*model.Comment, error) {
	var commet model.Comment
	query := "SELECT * from comments WHERE id = $1"
	err := cs.db.QueryRow(query, id).Scan(&commet.ID, &commet.PostID, &commet.UserID, &commet.Content, &commet.Image, &commet.CreatedAt, &commet.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &commet, nil
}

func (cs CommentStore) CreateComment(comment model.Comment) error {
	var query string
	if comment.Image.Valid {
		query = "INSERT INTO comments (id,post_id, user_id, content, image) VALUES ($1, $2, $3, $4, $5)"
		_, err := cs.db.Exec(query, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.Image)
		if err != nil {
			return err
		}

	} else {
		query = "INSERT INTO comments (id,post_id, user_id, content) VALUES ($1, $2, $3,$4)"
		_, err := cs.db.Exec(query, comment.ID, comment.PostID, comment.UserID, comment.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs CommentStore) UpdateComment(comment model.Comment) error {
	query := "UPDATE comments SET content = $1, image = $2 WHERE id = $3"
	_, err := cs.db.Exec(query, comment.Content, comment.Image, comment.ID)
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
