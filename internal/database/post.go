package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BasePostStore interface {
	GetPosts(pagination Pagination, search Search) ([]model.Post, error)
	GetPostByID(id int64) (*model.Post, error)
	GetPostsByUserID(userID int64, pagination Pagination, search Search) ([]model.Post, error)
	CreatePost(post model.Post) error
	UpdatePost(post model.Post) error
	DeletePost(id int64) error
}

type PostStore struct {
	DB *sql.DB
}

func NewPostStore(db *sql.DB) BasePostStore {
	return &PostStore{DB: db}
}

func (s *PostStore) GetPosts(pagination Pagination, search Search) ([]model.Post, error) {
	var posts []model.Post
	query := `
	WITH limited_posts AS (
		SELECT * FROM posts
		WHERE content ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	)
	SELECT posts.id,posts.content,posts.user_id,posts.image,posts.created_at,posts.updated_at,
    post_user.id,post_user.name,post_user.last_name,post_user.username,post_user.email,
    comments.id,comments.post_id,comments.user_id,comments.content,comments.image,
    comment_user.id,comment_user.name,comment_user.last_name,comment_user.username,comment_user.email
    FROM limited_posts as posts 
    JOIN users as post_user ON posts.user_id = post_user.id
    LEFT JOIN comments ON posts.id = comments.post_id
    LEFT JOIN users as comment_user ON comments.user_id = comment_user.id`

	rows, err := s.DB.Query(query, search.Query, pagination.Limit, pagination.Offset)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	postMap := make(map[int64]*model.Post)
	for rows.Next() {
		post := model.Post{}
		comment := model.Comment{}
		var commentID uuid.UUID
		var commentPostID, commentUserID, commentUsersID sql.NullInt64
		var commentContent, commentImage sql.NullString
		var commentUserName, commentUserLastName, commentUserUsername, commentUserEmail sql.NullString

		err := rows.Scan(&post.ID, &post.Content, &post.UserID, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username, &post.User.Email,
			&commentID, &commentPostID, &commentUserID, &commentContent, &commentImage,
			&commentUsersID, &commentUserName, &commentUserLastName, &commentUserUsername, &commentUserEmail)
		if err != nil {
			return nil, err
		}
		if _, ok := postMap[post.ID]; !ok {
			postMap[post.ID] = &post
		}
		if commentID != uuid.Nil {
			comment.ID = commentID
			comment.PostID = commentPostID.Int64
			comment.UserID = commentUserID.Int64
			comment.Content = commentContent.String
			comment.Image = sql.NullString{String: commentImage.String, Valid: commentImage.Valid}
			comment.User.ID = commentUsersID.Int64
			comment.User.Name = commentUserName.String
			comment.User.LastName = commentUserLastName.String
			comment.User.Username = commentUserUsername.String
			comment.User.Email = commentUserEmail.String

			postMap[comment.PostID].Comments = append(postMap[comment.PostID].Comments, comment)
		}

	}
	if err := rows.Err(); err != nil {

		return nil, err
	}
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	return posts, nil
}

func (s *PostStore) GetPostByID(id int64) (*model.Post, error) {
	post := &model.Post{}

	// First query: Get post details
	postQuery := `
		SELECT posts.id, posts.content, posts.user_id, posts.image, posts.created_at, posts.updated_at,
        users.id, users.name, users.last_name, users.username, users.email
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = $1`

	err := s.DB.QueryRow(postQuery, id).Scan(
		&post.ID, &post.Content, &post.UserID, &post.Image, &post.CreatedAt, &post.UpdatedAt,
		&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username, &post.User.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	commentQuery := `SELECT c.id, c.post_id, c.user_id, c.content, c.image,
        u.id, u.name, u.last_name, u.username, u.email
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = $1`

	commentRows, err := s.DB.Query(commentQuery, id)
	if err != nil {
		return nil, err
	}
	defer commentRows.Close()

	for commentRows.Next() {
		comment := model.Comment{}
		err := commentRows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image,
			&comment.User.ID, &comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email)
		if err != nil {
			return nil, err
		}
		post.Comments = append(post.Comments, comment)
	}

	return post, nil

}

func (s *PostStore) GetPostsByUserID(userID int64, pagination Pagination, search Search) ([]model.Post, error) {
	posts := []model.Post{}

	postQuery := `
	WITH limited_posts AS (
		SELECT * FROM posts
		WHERE user_id = $1
		AND content ILIKE '%' || $2 || '%'
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	)
	SELECT posts.id, posts.content, posts.user_id, posts.image, posts.created_at, posts.updated_at,
        users.id, users.name, users.last_name, users.username, users.email,
		comment_user.id, comment_user.name, comment_user.last_name, comment_user.username, comment_user.email
        FROM limited_posts as posts
        JOIN users ON posts.user_id = users.id
		LEFT JOIN comments ON posts.id = comments.post_id
		LEFT JOIN users as comment_user ON comments.user_id = comment_user.id`

	rows, err := s.DB.Query(postQuery, userID, search.Query, pagination.Limit, pagination.Offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	postMap := make(map[int64]*model.Post)
	for rows.Next() {
		post := model.Post{}
		comment := model.Comment{}
		err := rows.Scan(&post.ID, &post.Content, &post.UserID, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username, &post.User.Email,
			&comment.User.ID, &comment.User.Name, &comment.User.LastName, &comment.User.Username, &comment.User.Email,
		)
		if err != nil {

			return nil, err
		}
		if _, ok := postMap[post.ID]; !ok {
			postMap[post.ID] = &post
		}
		if comment.User.ID != 0 {
			postMap[post.ID].Comments = append(postMap[post.ID].Comments, comment)
		}
	}
	if err := rows.Err(); err != nil {

		return nil, err
	}

	for _, post := range postMap {
		posts = append(posts, *post)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}

func (s *PostStore) CreatePost(post model.Post) error {
	var query string

	if post.Image.Valid {
		query = "INSERT INTO posts (content, user_id, image) VALUES ($1, $2, $3)"

		_, err := s.DB.Exec(query, post.Content, post.UserID, post.Image)
		if err != nil {
			return err
		}

	} else {
		query = "INSERT INTO posts (content, user_id) VALUES ($1, $2)"

		_, err := s.DB.Exec(query, post.Content, post.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostStore) UpdatePost(post model.Post) error {
	if post.Image.Valid {
		query := "UPDATE posts SET content = $1, image = $2 WHERE id = $3"
		_, err := s.DB.Exec(query, post.Content, post.Image, post.ID)
		if err != nil {
			return err
		}
	} else {
		query := "UPDATE posts SET content = $1 WHERE id = $2"
		_, err := s.DB.Exec(query, post.Content, post.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostStore) DeletePost(id int64) error {
	query := "DELETE FROM posts WHERE id = $1"
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
