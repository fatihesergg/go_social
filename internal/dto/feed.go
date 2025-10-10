package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type FeedResponse struct {
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

func NewFeedResponse(posts []model.Post) []FeedResponse {
	result := []FeedResponse{}
	for _, post := range posts {
		feedResponse := FeedResponse{
			ID:           post.ID,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			User:         post.User,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			IsLiked:      post.IsLiked,
			IsFollowing:  post.IsFollowing,
		}
		result = append(result, feedResponse)
	}
	return result
}
