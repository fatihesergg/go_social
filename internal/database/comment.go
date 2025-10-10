package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentStore interface {
	GetCommentsByPostID(postID, userID uuid.UUID) ([]model.Comment, error)
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

func (cs CommentStore) GetCommentsByPostID(postID, userID uuid.UUID) ([]model.Comment, error) {
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
	(user_follows.follow_id IS NOT NULL) AS is_following,

	replies.id,
	replies.comment_id,
	replies.message,

	reply_user.id,
	reply_user.name,
	reply_user.last_name,
	reply_user.username



	FROM comments JOIN users ON comments.user_id = users.id
	LEFT JOIN comment_likes_count ON comment_likes_count.comment_id = comments.id
	LEFT JOIN reply_count ON reply_count.comment_id = comments.id
	LEFT JOIN user_likes ON user_likes.comment_id = comments.id
	LEFT JOIN user_follows ON user_follows.follow_id = users.id
	LEFT JOIN replies ON replies.comment_id = comments.id
	LEFT JOIN users AS reply_user ON reply_user.id = users.id
	WHERE comments.post_id = $1`
	rows, err := cs.db.Query(query, postID, userID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	commentMap := make(map[uuid.UUID]*model.Comment)
	for rows.Next() {
		var comment model.Comment
		var replyID, replyUserID, replyCommentID *uuid.UUID
		var replyMessage *string
		var replyUserName, replyUserLastName, replyUserUsername *string
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.ID, &comment.User.Name, &comment.User.LastName, &comment.User.Username,
			&comment.LikeCount, &comment.ReplyCount,
			&comment.IsLiked, &comment.IsFollowing,
			&replyID, &replyCommentID, &replyMessage,
			&replyUserID, &replyUserName, &replyUserLastName, &replyUserUsername,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if _, ok := commentMap[comment.ID]; !ok {
			commentMap[comment.ID] = &comment
		}
		if replyID != nil {
			reply := model.Reply{
				ID:        *replyID,
				CommentID: *replyCommentID,
				Message:   *replyMessage,
				User: model.User{
					ID:       *replyUserID,
					Name:     *replyUserName,
					LastName: *replyUserLastName,
					Username: *replyUserName,
				},
			}
			commentMap[*replyCommentID].Replies = append(commentMap[*replyCommentID].Replies, reply)
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
	query := `SELECT comments.id,comments.post_id,comments.user_id,comments.content,comments.created_at,comments.updated_at,
		users.name,users.last_name,users.username,users.email	
	FROM comments 
	JOIN users ON comments.user_id = users.id
	WHERE comments.id = $1`
	err := cs.db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if comment.ID == uuid.Nil {
		return nil, nil
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
