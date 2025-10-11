package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseReplyStore interface {
	CreateReply(reply *model.Reply) error
	UpdateReply(reply *model.Reply) error
	GetRepliesByCommentID(commentID uuid.UUID) ([]model.Reply, error)
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
	reply := model.Reply{}
	query := "SELECT id,comment_id,user_id,message FROM replies WHERE id = $1"
	err := rs.DB.QueryRow(query, replyID).Scan(&reply.ID, &reply.CommentID, &reply.UserID, &reply.Message)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reply, err

}

func (rc *ReplyStore) DeleteReply(replyID uuid.UUID) error {
	query := "DELETE FROM replies WHERE id = $1"
	_, err := rc.DB.Exec(query, replyID)
	return err
}
func (rc *ReplyStore) GetRepliesByCommentID(commentID uuid.UUID) ([]model.Reply, error) {
	replies := []model.Reply{}
	query := `

	SELECT

	replies.id,
	replies.message,

	reply_user.id,
	reply_user.name,
	reply_user.last_name,
	reply_user.username

	FROM replies
	LEFT JOIN users as reply_user ON reply_user.id = replies.user_id
	WHERE replies.comment_id = $1
	`
	rows, err := rc.DB.Query(query, commentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		reply := model.Reply{}

		err := rows.Scan(
			&reply.ID, &reply.Message,
			&reply.User.ID, &reply.User.Name, &reply.User.LastName, &reply.User.Username,
		)
		if err != nil {
			return nil, err
		}

		replies = append(replies, reply)

	}
	if rows.Err() != nil {
		return nil, err
	}
	if len(replies) == 0 {
		return nil, nil
	}
	return replies, nil

}
