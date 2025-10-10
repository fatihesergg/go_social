package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
)

type BaseReplyStore interface {
	CreateReply(reply *model.Reply) error
	UpdateReply(reply *model.Reply) error
}

type ReplyStore struct {
	DB *sql.DB
}

func NewReplyStore(db *sql.DB) *ReplyStore {
	return &ReplyStore{
		DB: db,
	}
}

func (rs *ReplyStore) CreateReply(reply *model.Reply) error {
	query := "INSERT INTO replies ( comment_id,user_id,message ) VALUES ( $1,$2,$3 )"
	_, err := rs.DB.Exec(query, reply.CommentID, reply.UserID, reply.Message)
	return err
}

func (rs *ReplyStore) UpdateReply(reply *model.Reply) error {
	query := "UPDATE replies SET comment_id = $1, user_id = $2, message = $3 )"
	_, err := rs.DB.Exec(query, reply.CommentID, reply.UserID, reply.Message)
	return err
}

//TODO: Get replies and get reply by id
