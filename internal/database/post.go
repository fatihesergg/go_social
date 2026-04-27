package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BasePostStore interface {
	HasAccessToPost(userID, postID uuid.UUID) (bool, error)
	GetPostByID(postID uuid.UUID) (*model.Post, error)
	GetPosts(pagination Pagination, search Search, userID uuid.UUID) ([]model.Post, error)
	GetPostsByTag(pagination Pagination, tag string, userID uuid.UUID) ([]model.Post, error)
	GetPostDetailsByID(postID, userID uuid.UUID) (*model.Post, error)
	GetPostsByUserID(userID uuid.UUID, pagination Pagination, search Search) ([]model.Post, error)
	CreatePost(post *model.Post) error
	UpdatePost(post *model.Post) error
	DeletePost(id uuid.UUID) error
}

type PostStore struct {
	DB *sql.DB
}

func NewPostStore(db *sql.DB) BasePostStore {
	return &PostStore{DB: db}
}

func (s *PostStore) HasAccessToPost(userID, postID uuid.UUID) (bool, error) {
	var result bool
	query := `
	SELECT EXISTS(SELECT 1 FROM posts
	LEFT JOIN follows ON follows.follow_id = posts.user_id
	WHERE posts.id = $2  
	AND ( 
	posts.visibility = 'public' 
	OR posts.user_id = $1
	OR ( follows.user_id IS NOT NULL AND follows.user_id = $1))
	)`

	err := s.DB.QueryRow(query, userID, postID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("Error while checking post access: %w", err)
	}
	return result, nil

}

func (s *PostStore) GetPosts(pagination Pagination, search Search, userID uuid.UUID) ([]model.Post, error) {
	var posts []model.Post

	query := `
	WITH likes_count AS (
		SELECT post_id ,COUNT(*) as total_likes FROM post_likes
		GROUP BY post_id
	),

	comments_count AS (
		SELECT post_id, COUNT(*) as total_comments FROM comments
		GROUP BY post_id
	),

	user_likes AS (
		SELECT post_id FROM post_likes
		WHERE user_id = $4
	),

	user_follows AS (
		SELECT follow_id
		FROM follows
		WHERE user_id = $5
	),

	limited_posts AS (
		SELECT * FROM posts
		WHERE content ILIKE '%' || $1 || '%' AND (posts.visibility = 'public' OR posts.user_id IN (SELECT follow_id FROM user_follows))
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	)



	SELECT 
	posts.id,posts.content,posts.created_at,posts.updated_at,
    post_user.id,post_user.name,post_user.last_name,post_user.username,
	
	COALESCE(likes_count.total_likes,0),
	COALESCE(comments_count.total_comments,0),

	(user_likes.post_id IS NOT NULL),
	(user_follows.follow_id IS NOT NULL)

    FROM limited_posts as posts 
    JOIN users as post_user ON posts.user_id = post_user.id
    LEFT JOIN likes_count ON likes_count.post_id = posts.id
	LEFT JOIN comments_count ON comments_count.post_id = posts.id
	LEFT JOIN user_likes ON user_likes.post_id = posts.id
	LEFT JOIN user_follows ON user_follows.follow_id = post_user.id
;
	`

	rows, err := s.DB.Query(query, search.Query, pagination.Limit, pagination.Offset, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		var commentCount, postLikeCount *int
		var isLiked, isFollowing *bool

		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&postLikeCount, &commentCount,
			&isLiked, &isFollowing,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		post.LikeCount = *postLikeCount
		post.CommentCount = *commentCount
		post.IsLiked = *isLiked
		post.IsFollowing = *isFollowing

		posts = append(posts, post)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}
func (s *PostStore) GetPostsByTag(pagination Pagination, tag string, userID uuid.UUID) ([]model.Post, error) {
	var posts []model.Post

	query := `
	WITH likes_count AS (
		SELECT post_id ,COUNT(*) as total_likes FROM post_likes
		GROUP BY post_id
	),

	comments_count AS (
		SELECT post_id, COUNT(*) as total_comments FROM comments
		GROUP BY post_id
	),

	user_likes AS (
		SELECT post_id FROM post_likes
		WHERE user_id = $4
	),

	user_follows AS (
		SELECT follow_id
		FROM follows
		WHERE user_id = $4
	),

	tag_post_ids AS (
	SELECT post_id FROM post_tags
	JOIN tags ON tags.id = post_tags.tag_id AND tags.name = $1
	),

	limited_posts AS (
		SELECT posts.* FROM tag_post_ids
		JOIN posts ON tag_post_ids.post_id = posts.id
		WHERE posts.visibility = 'public' OR posts.user_id IN (SELECT follow_id FROM user_follows)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	)



	SELECT 
	posts.id,posts.content,posts.created_at,posts.updated_at,
    post_user.id,post_user.name,post_user.last_name,post_user.username,
	
	COALESCE(likes_count.total_likes,0),
	COALESCE(comments_count.total_comments,0),

	(user_likes.post_id IS NOT NULL),
	(user_follows.follow_id IS NOT NULL)

    FROM limited_posts as posts 
    JOIN users as post_user ON posts.user_id = post_user.id
    LEFT JOIN likes_count ON likes_count.post_id = posts.id
	LEFT JOIN comments_count ON comments_count.post_id = posts.id
	LEFT JOIN user_likes ON user_likes.post_id = posts.id
	LEFT JOIN user_follows ON user_follows.follow_id = post_user.id
;
	`

	rows, err := s.DB.Query(query, tag, pagination.Limit, pagination.Offset, userID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all posts by tags: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		var commentCount, postLikeCount *int
		var isLiked, isFollowing *bool

		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&postLikeCount, &commentCount,
			&isLiked, &isFollowing,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		post.LikeCount = *postLikeCount
		post.CommentCount = *commentCount
		post.IsLiked = *isLiked
		post.IsFollowing = *isFollowing

		posts = append(posts, post)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}

func (s *PostStore) GetPostByID(postID uuid.UUID) (*model.Post, error) {
	result := &model.Post{}
	var id, userID *uuid.UUID
	var content *string
	query := `SELECT id,user_id,content FROM posts WHERE id = $1`
	err := s.DB.QueryRow(query, postID).Scan(&id, &userID, &content)
	if err != nil {
		return nil, fmt.Errorf("Error while getting post by id: %w", err)
	}
	result.ID = *id
	result.Content = *content
	result.UserID = *userID

	return result, nil
}

func (s *PostStore) GetPostDetailsByID(postID, userID uuid.UUID) (*model.Post, error) {

	postQuery := `


		WITH reply_count AS (
			SELECT comment_id,COUNT(*) as total_reply FROM replies
			GROUP BY comment_id
		),

		comment_count AS (
			SELECT post_id,COUNT(*) AS total_comment  FROM comments
			GROUP BY post_id
		),

		comment_like_count AS (
			SELECT comment_id,COUNT(*) AS total_comment_like FROM comment_likes
			GROUP BY comment_id
		),

		post_like_count AS (
			SELECT post_id,COUNT(*) AS total_post_like FROM post_likes
			GROUP BY post_id
		),

		user_follows AS (
			SELECT follow_id FROM follows
			WHERE user_id = $1
		),

		user_comment_likes AS  (
			SELECT comment_id FROM comment_likes
			WHERE user_id  = $1
		),

		user_post_likes AS  (
			SELECT post_id FROM post_likes
			WHERE user_id  = $1
		)

		SELECT 
		posts.id,
		posts.content,
		posts.created_at,
		posts.updated_at,
		posts.visibility,

        post_user.id,
		post_user.name,
		post_user.last_name,
		post_user.username,

		comments.id,
		comments.content,

		comment_user.id,
		comment_user.name,
		comment_user.last_name,
		comment_user.username,

		COALESCE(comment_count.total_comment,0) AS  total_comment,
		COALESCE(reply_count.total_reply,0) AS total_reply,
		COALESCE(comment_like_count.total_comment_like,0) AS total_comment_like,
		COALESCE(post_like_count.total_post_like,0) AS total_post_like,
		
		(user_post_likes.post_id IS NOT NULL) AS is_post_liked,
		(user_comment_likes.comment_id IS NOT NULL) is_comment_liked,

		(post_follows.follow_id IS NOT NULL) AS is_post_following,
		(comment_follows.follow_id IS NOT NULL) AS is_comment_following

        FROM posts
        JOIN users AS post_user ON posts.user_id = post_user.id
		LEFT JOIN comments ON comments.post_id = posts.id
		LEFT JOIN users AS comment_user ON comment_user.id = comments.user_id
		LEFT JOIN user_follows AS post_follows ON post_follows.follow_id = post_user.id
		LEFT JOIN user_follows AS comment_follows ON comment_follows.follow_id = comments.user_id
		LEFT JOIN user_post_likes ON user_post_likes.post_id = posts.id
		LEFT JOIN user_comment_likes ON user_comment_likes.comment_id  = comments.id
		LEFT JOIN comment_count ON  comment_count.post_id  = posts.id
		LEFT JOIN reply_count ON  reply_count.comment_id = comments.id
		LEFT JOIN comment_like_count ON  comment_like_count.comment_id = comments.id
		LEFT JOIN post_like_count ON  post_like_count.post_id = posts.id
        WHERE posts.id = $2`

	rows, err := s.DB.Query(postQuery, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting post details by id: %w", err)
	}

	defer rows.Close()

	post := &model.Post{}
	for rows.Next() {
		var commentID, commentUserID *uuid.UUID
		var commentContent *string
		var commentUserName *string
		var commentUserLastName *string
		var commentUserUsername *string
		var replyCount, commentLikeCount *int
		var isCommentFollowing, isCommentLiked *bool

		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Visibility,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&commentID, &commentContent,
			&commentUserID, &commentUserName, &commentUserLastName, &commentUserUsername,
			&post.CommentCount, &replyCount, &commentLikeCount, &post.LikeCount,
			&post.IsLiked, &isCommentLiked,
			&post.IsFollowing, &isCommentFollowing,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		if commentID != nil {

			comment := &model.Comment{
				ID:      *commentID,
				Content: *commentContent,
				User: model.User{
					ID:       *commentUserID,
					Name:     *commentUserName,
					LastName: *commentUserLastName,
					Username: *commentUserUsername,
				},
				IsLiked:     *isCommentLiked,
				IsFollowing: *isCommentFollowing,
				ReplyCount:  *replyCount,
				LikeCount:   *commentLikeCount,
			}
			post.Comments = append(post.Comments, *comment)

		}

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}
	if post.ID == uuid.Nil {
		return nil, sql.ErrNoRows
	}

	return post, nil

}

func (s *PostStore) GetPostsByUserID(userID uuid.UUID, pagination Pagination, search Search) ([]model.Post, error) {
	posts := []model.Post{}

	postQuery := `
	WITH limited_posts AS (
		SELECT * FROM posts
		WHERE user_id = $1
		AND content ILIKE '%' || $2 || '%'
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	)
	SELECT posts.id, posts.content, posts.created_at, posts.updated_at,
        users.id,users.name, users.last_name, users.username,
		comments.id,comments.content,comment_user.name, comment_user.last_name, comment_user.username
        FROM limited_posts as posts
        JOIN users ON posts.user_id = users.id
		LEFT JOIN comments ON posts.id = comments.post_id
		LEFT JOIN users as comment_user ON comments.user_id = comment_user.id`

	rows, err := s.DB.Query(postQuery, userID.String(), search.Query, pagination.Limit, pagination.Offset)

	if err != nil {
		return nil, fmt.Errorf("Error while getting posts by user id: %w", err)
	}
	defer rows.Close()
	postMap := make(map[uuid.UUID]*model.Post)
	for rows.Next() {
		post := model.Post{}
		var commentID *uuid.UUID
		var commentContent *string
		var commentUserName *string
		var commentUserLastName *string
		var commentUserUsername *string
		err := rows.Scan(&post.ID, &post.Content, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username, &commentID, &commentContent,
			&commentUserName, &commentUserLastName, &commentUserUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning result: %w", err)
		}
		if _, ok := postMap[post.ID]; !ok {
			postMap[post.ID] = &post
		}
		if commentID != nil {
			comment := model.Comment{
				Content: *commentContent,
				User: model.User{
					Name:     *commentUserName,
					LastName: *commentUserLastName,
					Username: *commentUserName,
				},
			}
			postMap[post.ID].Comments = append(postMap[post.ID].Comments, comment)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in result row: %w", err)
	}

	for _, post := range postMap {
		posts = append(posts, *post)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}

func (s *PostStore) CreatePost(post *model.Post) error {

	query := "INSERT INTO posts (id, content, user_id, visibility) VALUES ($1, $2, $3,$4)"

	_, err := s.DB.Exec(query, post.ID, post.Content, post.UserID, post.Visibility)
	if err != nil {
		return fmt.Errorf("Error while inserting post: %w", err)
	}

	return nil
}

func (s *PostStore) UpdatePost(post *model.Post) error {
	query := "UPDATE posts SET content = $1 WHERE id = $2"
	result, err := s.DB.Exec(query, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("Error while updating post: %w", err)
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

func (s *PostStore) DeletePost(id uuid.UUID) error {
	query := "DELETE FROM posts WHERE id = $1"
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error while deleting post: %w", err)
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
