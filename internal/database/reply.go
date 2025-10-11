package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseReplyStore interface {
	CreateReply(reply *model.Reply) error
	UpdateReply(reply *model.Reply) error
	GetReplyByID(replyID uuid.UUID) (*model.Reply, error)
	DeleteReply(replyID uuid.UUID) error
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
	query := "UPDATE replies SET comment_id = $1, user_id = $2, message = $3 WHERE id = $4"
	_, err := rs.DB.Exec(query, reply.CommentID, reply.UserID, reply.Message, reply.ID)
	return err
}

func (rs *ReplyStore) GetReplyByID(replyID uuid.UUID) (*model.Reply, error) {
	var replyid, replyCommentID, replyUserID *uuid.UUID
	var replyMessage *string
	query := "SELECT id,comment_id,user_id,message FROM replies WHERE id = $1"
	err := rs.DB.QueryRow(query, replyID).Scan(&replyid, &replyCommentID, &replyUserID, &replyMessage)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model.Reply{ID: *replyid, CommentID: *replyCommentID, UserID: *replyUserID, Message: *replyMessage}, err

}

func (rc *ReplyStore) DeleteReply(replyID uuid.UUID) error {
	query := "DELETE FROM replies WHERE id = $1"
	_, err := rc.DB.Exec(query, replyID)
	return err
}
