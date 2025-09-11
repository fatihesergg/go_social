package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/lib/pq"
)

type BaseFeedStore interface {
	GetFeed(userID int64) ([]model.Post, error)
}

type FeedStore struct {
	DB *sql.DB
}

func NewFeedStore(db *sql.DB) BaseFeedStore {
	return &FeedStore{
		DB: db,
	}
}

func (fs FeedStore) GetFeed(userID int64) ([]model.Post, error) {
	var posts []model.Post

	followers := []int64{}
	followersQuery := `SELECT follow_id FROM follows WHERE user_id = $1`
	followersRows, err := fs.DB.Query(followersQuery, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer followersRows.Close()
	for followersRows.Next() {
		var followerID int64
		if err := followersRows.Scan(&followerID); err != nil {
			fmt.Println(err)
			return nil, err
		}
		followers = append(followers, followerID)
	}

	query := `SELECT posts.id, posts.content, posts.image, posts.created_at, posts.updated_at,
	 post_user.name, post_user.last_name, post_user.username,
	comments.id, comments.content, comments.post_id, comments.created_at, comments.updated_at,
	comment_user.name, comment_user.last_name, comment_user.username
	FROM posts
	JOIN users as post_user ON posts.user_id = post_user.id
	LEFT JOIN comments ON posts.id = comments.post_id
	JOIN users as comment_user ON comments.user_id = comment_user.id
	WHERE posts.user_id = ANY($1)`

	rows, err := fs.DB.Query(query, pq.Array(followers))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	postMap := make(map[int64]*model.Post)
	for rows.Next() {
		post := model.Post{}
		comment := model.Comment{}
		comment.User = model.User{}
		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.Name, &post.User.LastName, &post.User.Username,
			&comment.ID, &comment.Content, &comment.PostID, &comment.CreatedAt, &comment.UpdatedAt,
			&comment.User.Name, &comment.User.LastName, &comment.User.Username)
		if err != nil {
			return nil, err
		}
		if _, ok := postMap[post.ID]; !ok {
			postMap[post.ID] = &post
		}
		postMap[comment.PostID].Comments = append(postMap[comment.PostID].Comments, comment)
	}

	for _, post := range postMap {
		posts = append(posts, *post)
	}

	return posts, nil
}
