package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseFeedStore interface {
	GetFeed(userID uuid.UUID, pagination Pagination, search Search) ([]model.Post, error)
}

type FeedStore struct {
	DB *sql.DB
}

func NewFeedStore(db *sql.DB) BaseFeedStore {
	return &FeedStore{
		DB: db,
	}
}

func (fs FeedStore) GetFeed(userID uuid.UUID, pagination Pagination, search Search) ([]model.Post, error) {
	var posts []model.Post

	query := `
	WITH limited_posts AS (
		SELECT * FROM posts
		WHERE content ILIKE '%' || $4 || '%' AND user_id = ANY (SELECT follow_id FROM follows WHERE user_id = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	),

	likes_count AS (
		SELECT post_id ,COUNT(*) as total_likes FROM post_likes
		GROUP BY post_id
	),

	comments_count AS (
		SELECT post_id, COUNT(*) as total_comments FROM comments
		GROUP BY post_id
	),

	user_likes AS (
		SELECT post_id FROM post_likes
		WHERE user_id = $1
	)


	SELECT 
	posts.id,
	posts.content,
	posts.created_at,
	posts.updated_at,

    post_user.id,
	post_user.name,
	post_user.last_name,
	post_user.username,
	
	COALESCE(likes_count.total_likes,0) AS total_likes,
	COALESCE(comments_count.total_comments,0) AS total_comments,

	(user_likes.post_id IS NOT NULL) AS is_liked

    FROM limited_posts as posts 
    JOIN users AS post_user ON post_user.id = posts.user_id
    LEFT JOIN likes_count ON likes_count.post_id = posts.id
	LEFT JOIN comments_count ON comments_count.post_id = posts.id
	LEFT JOIN user_likes ON user_likes.post_id = posts.id
	`

	rows, err := fs.DB.Query(query, userID, pagination.Limit, pagination.Offset, search.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&post.LikeCount, &post.CommentCount,
			&post.IsLiked,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}
	if err := rows.Err(); err != nil {

		return nil, err
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}
