package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type CreatePostDTO struct {
	Content string `json:"content" binding:"required,lte=500"`
	Image   string `json:"image"`
}

type UpdatePostDTO struct {
	Content string `json:"content" binding:"required,alphanum,lte=500"`
	Image   string `json:"image"`
}

type AllPostRepsonse struct {
	ID           uuid.UUID  `json:"id"`
	Content      string     `json:"content"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	User         model.User `json:"user"`
	LikeCount    int        `json:"total_likes"`
	CommentCount int        `json:"total_comment"`
	IsLiked      bool       `json:"is_liked"`
	IsFollowing  bool       `json:"is_following"`
}

type PostDetailResponse struct {
	ID           uuid.UUID         `json:"id"`
	Content      string            `json:"content"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	User         model.User        `json:"user"`
	Comments     []CommentResponse `json:"comments"`
	LikeCount    int               `json:"total_likes"`
	CommentCount int               `json:"total_comment"`
	IsLiked      bool              `json:"is_liked"`
	IsFollowing  bool              `json:"is_following"`
}

func NewAllPostResponse(posts []model.Post) []AllPostRepsonse {
	result := []AllPostRepsonse{}
	for _, post := range posts {
		result = append(result, AllPostRepsonse{
			ID:           post.ID,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			User:         post.User,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			IsLiked:      post.IsLiked,
			IsFollowing:  post.IsFollowing,
		})
	}
	return result
}

func NewPostDetailResponse(post *model.Post) PostDetailResponse {
	result := PostDetailResponse{
		ID:           post.ID,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
		User:         post.User,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		IsLiked:      post.IsLiked,
		IsFollowing:  post.IsFollowing,
		Comments:     NewCommentResponse(post.Comments),
	}
	return result
}
