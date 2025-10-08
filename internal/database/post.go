package database

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BasePostStore interface {
	GetPosts(pagination Pagination, search Search) ([]model.Post, error)
	GetPostByID(id uuid.UUID) (*model.Post, error)
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

func (s *PostStore) GetPosts(pagination Pagination, search Search) ([]model.Post, error) {
	var posts []model.Post
	query := `
	WITH limited_posts AS (
		SELECT * FROM posts
		WHERE content ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	)
	SELECT posts.id,posts.content,posts.image,posts.created_at,posts.updated_at,
    post_user.id,post_user.name,post_user.last_name,post_user.username,
    comments.id,comments.post_id,comments.content,comments.image,
    comment_user.name,comment_user.last_name,comment_user.username
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
	postMap := make(map[uuid.UUID]*model.Post)
	for rows.Next() {
		post := model.Post{}
		var commentID, commentPostID *uuid.UUID
		var commentContent *string
		var commentImage *string
		var commentUserName *string
		var commentUserLastName *string
		var commentUserUsername *string

		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&commentID, &commentPostID, &commentContent, &commentImage,
			&commentUserName, &commentUserLastName, &commentUserUsername,
		)
		if err != nil {
			return nil, err
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
					Username: *commentUserUsername,
				},
			}
			postMap[post.ID].Comments = append(postMap[post.ID].Comments, comment)
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

func (s *PostStore) GetPostByID(id uuid.UUID) (*model.Post, error) {

	postQuery := `
		SELECT posts.id, posts.content,  posts.image, posts.created_at, posts.updated_at,
        users.id,users.name, users.last_name, users.username,
		comments.id,comments.post_id,comments.content,
		comment_user.name,comment_user.last_name,comment_user.username
        FROM posts
        JOIN users ON posts.user_id = users.id
		LEFT JOIN comments ON comments.post_id = posts.id
		LEFT JOIN users as comment_user ON comment_user.id = comments.user_id
        WHERE posts.id = $1`

	rows, err := s.DB.Query(postQuery, id.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	post := &model.Post{}
	for rows.Next() {
		var commentID, commentPostID *uuid.UUID
		var commentContent *string
		var commentUserName *string
		var commentUserLastName *string
		var commentUserUsername *string

		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username,
			&commentID, &commentPostID, &commentContent,
			&commentUserName, &commentUserLastName, &commentUserUsername,
		)
		if err != nil {
			return nil, err
		}
		if commentID != nil {
			comment := model.Comment{
				Content: *commentContent,
				User: model.User{
					Name:     *commentUserName,
					LastName: *commentUserLastName,
					Username: *commentUserUsername,
				},
			}
			post.Comments = append(post.Comments, comment)
		}

	}
	if err := rows.Err(); err != nil {

		return nil, err
	}
	if post.ID == uuid.Nil {
		return nil, nil
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
	SELECT posts.id, posts.content,posts.image, posts.created_at, posts.updated_at,
        users.id,users.name, users.last_name, users.username,
		comments.id,comments.content,comment_user.name, comment_user.last_name, comment_user.username
        FROM limited_posts as posts
        JOIN users ON posts.user_id = users.id
		LEFT JOIN comments ON posts.id = comments.post_id
		LEFT JOIN users as comment_user ON comments.user_id = comment_user.id`

	rows, err := s.DB.Query(postQuery, userID.String(), search.Query, pagination.Limit, pagination.Offset)

	if err != nil {
		return nil, err
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
		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.CreatedAt, &post.UpdatedAt,
			&post.User.ID, &post.User.Name, &post.User.LastName, &post.User.Username, &commentID, &commentContent,
			&commentUserName, &commentUserLastName, &commentUserUsername,
		)
		if err != nil {

			return nil, err
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

func (s *PostStore) CreatePost(post *model.Post) error {

	query := "INSERT INTO posts (content, user_id) VALUES ($1, $2)"

	_, err := s.DB.Exec(query, post.Content, post.UserID.String())
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) UpdatePost(post *model.Post) error {
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

func (s *PostStore) DeletePost(id uuid.UUID) error {
	query := "DELETE FROM posts WHERE id = $1"
	_, err := s.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
