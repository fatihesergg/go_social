package database

import (
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
)

type BasePostStore interface {
	GetPosts() ([]model.Post, error)
	GetPostByID(id int64) (*model.Post, error)
	GetPostsByUserID(userID int64) ([]model.Post, error)
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

func (s *PostStore) GetPosts() ([]model.Post, error) {
	var posts []model.Post
	query := "SELECT * FROM posts"
	rows, err := s.DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		post := model.Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostStore) GetPostByID(id int64) (*model.Post, error) {
	post := &model.Post{}
	query := "SELECT * FROM posts WHERE id = $1"
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&post.ID, &post.Content, &post.Image, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return post, nil
}

func (s *PostStore) GetPostsByUserID(userID int64) ([]model.Post, error) {
	posts := []model.Post{}
	query := "SELECT * FROM posts WHERE user_id = $1"
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := model.Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostStore) CreatePost(post model.Post) error {
	var query string
	fmt.Println(post)
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
