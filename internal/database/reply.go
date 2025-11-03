package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseReplyStore interface {
	HasAccessToReply(userID, replyID uuid.UUID) (bool, error)
	CreateCommentReply(reply *model.Reply) error
	CreateNestedReply(reply *model.Reply) error
	UpdateReply(reply *model.Reply) error
	GetRepliesByCommentID(commentID, userID uuid.UUID) ([]model.Reply, error)
	GetRepliesByParentID(parentID, userID uuid.UUID) ([]model.Reply, error)
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

func (rs *ReplyStore) HasAccessToReply(userID, replyID uuid.UUID) (bool, error) {
	var result bool
	query := `SELECT EXISTS(
	SELECT 1 FROM replies
	JOIN comments ON  replies.comment_id = comments.id
	JOIN posts ON posts.id = comments.post_id
	LEFT JOIN follows ON follows.following_id
	WHERE 
	replies.id = $2 
	AND (
	posts.visibility  = 'public' 
	OR posts.user_id = $1
	OR (follows.user_id IS NOT NULL AND follows.user_id = $1)
	)
	)`
	err := rs.DB.QueryRow(query, userID, replyID).Scan(&result)
	if err != nil {
		return false, err
	}
	return result, err

}

func (rs *ReplyStore) CreateCommentReply(reply *model.Reply) error {
	query := "INSERT INTO replies ( id,comment_id,user_id,message ) VALUES ( $1,$2,$3,$4 )"
	_, err := rs.DB.Exec(query, reply.ID, reply.CommentID, reply.UserID, reply.Message)
	return err
}

func (rs *ReplyStore) CreateNestedReply(reply *model.Reply) error {
	query := "INSERT INTO replies ( id,parent_id,user_id,message ) VALUES ( $1,$2,$3,$4 )"
	_, err := rs.DB.Exec(query, reply.ID, reply.ParentID, reply.UserID, reply.Message)
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
func (rc *ReplyStore) GetRepliesByCommentID(commentID, userID uuid.UUID) ([]model.Reply, error) {
	replies := []model.Reply{}
	query := `

	WITH reply_count AS (
		SELECT parent_id,COUNT(*) AS replies_count FROM replies
		GROUP BY parent_id
	),

	reply_likes_count AS (
	SELECT reply_id,COUNT(*) AS likes_count FROM reply_likes
	GROUP BY reply_id
	),

	user_likes AS (
	SELECT reply_id FROM reply_likes
	WHERE user_id = $2 
	),

	user_follows AS (
	SELECT follow_id FROM follows
	WHERE user_id = $2
	)

	SELECT

	replies.id,
	replies.message,

	reply_user.id,
	reply_user.name,
	reply_user.last_name,
	reply_user.username,

	COALESCE(reply_likes_count.likes_count ,0) AS total_likes,
	COALESCE(reply_count.replies_count ,0) AS total_replies,

	(user_likes.reply_id IS NOT NULL) AS is_liked,
	(user_follows.follow_id IS NOT NULL) AS is_following


	FROM replies
	JOIN users as reply_user ON reply_user.id = replies.user_id
	LEFT JOIN reply_count ON reply_count.parent_id = replies.id
	LEFT JOIN reply_likes_count ON reply_likes_count.reply_id = replies.id
	LEFT JOIN user_likes ON user_likes.reply_id = replies.id
	LEFT JOIN user_follows ON user_follows.follow_id = reply_user.id
	WHERE replies.comment_id = $1
	`
	rows, err := rc.DB.Query(query, commentID, userID)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		reply := model.Reply{}

		err := rows.Scan(
			&reply.ID, &reply.Message,
			&reply.User.ID, &reply.User.Name, &reply.User.LastName, &reply.User.Username,
			&reply.LikeCount, &reply.ReplyCount, &reply.IsLiked, &reply.IsFollowing,
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

func (rc *ReplyStore) GetRepliesByParentID(parentID, userID uuid.UUID) ([]model.Reply, error) {
	replies := []model.Reply{}
	query := `

	WITH reply_count AS (
		SELECT parent_id,COUNT(*) AS replies_count FROM replies
		GROUP BY parent_id
	),

	reply_likes_count AS (
	SELECT reply_id,COUNT(*) AS likes_count FROM reply_likes
	GROUP BY reply_id
	),

	user_likes AS (
	SELECT reply_id FROM reply_likes
	WHERE user_id = $2 
	),

	user_follows AS (
	SELECT follow_id FROM follows
	WHERE user_id = $2
	)

	SELECT

	replies.id,
	replies.parent_id,
	replies.message,

	reply_user.id,
	reply_user.name,
	reply_user.last_name,
	reply_user.username,

	COALESCE(reply_likes_count.likes_count ,0) AS total_likes,
	COALESCE(reply_count.replies_count ,0) AS total_replies,

	(user_likes.reply_id IS NOT NULL) AS is_liked,
	(user_follows.follow_id IS NOT NULL) AS is_following

	FROM replies
	JOIN users as reply_user ON reply_user.id = replies.user_id
	LEFT JOIN reply_count ON reply_count.parent_id = replies.id
	LEFT JOIN reply_likes_count ON reply_likes_count.reply_id = replies.id
	LEFT JOIN user_likes ON user_likes.reply_id = replies.id
	LEFT JOIN user_follows ON user_follows.follow_id = reply_user.id
	WHERE replies.parent_id = $1
	`
	rows, err := rc.DB.Query(query, parentID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		reply := model.Reply{}

		err := rows.Scan(
			&reply.ID, &reply.ParentID, &reply.Message,
			&reply.User.ID, &reply.User.Name, &reply.User.LastName, &reply.User.Username,
			&reply.LikeCount, &reply.ReplyCount, &reply.IsLiked, &reply.IsFollowing,
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
