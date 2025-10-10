package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type CreateCommentDTO struct {
	PostID  uuid.UUID `json:"post_id" binding:"required,uuid"`
	Content string    `json:"content" binding:"required,lte=200"`
	Image   string    `json:"image"`
}

type UpdateCommentDTO struct {
	Content string `json:"content" binding:"required,alphanum,lte=200"`
	Image   string `json:"image"`
}

type CommentResponse struct {
	ID          uuid.UUID  `json:"id"`
	Content     string     `json:"content"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	User        model.User `json:"user"`
	LikeCount   int        `json:"total_likes"`
	ReplyCount  int        `json:"total_reply"`
	IsLiked     bool       `json:"is_liked"`
	IsFollowing bool       `json:"is_followed"`
}

type CommentDetailResponse struct {
	ID          uuid.UUID       `json:"id"`
	Content     string          `json:"content"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	User        model.User      `json:"user"`
	Replies     []ReplyResponse `json:"replies"`
	LikeCount   int             `json:"total_likes"`
	ReplyCount  int             `json:"total_reply"`
	IsLiked     bool            `json:"is_liked"`
	IsFollowing bool            `json:"is_followed"`
}

func NewCommentResponse(comments []model.Comment) []CommentResponse {
	result := []CommentResponse{}
	for _, comment := range comments {
		commentResponse := CommentResponse{
			ID:          comment.ID,
			Content:     comment.Content,
			CreatedAt:   comment.CreatedAt,
			UpdatedAt:   comment.UpdatedAt,
			User:        comment.User,
			LikeCount:   comment.LikeCount,
			ReplyCount:  comment.ReplyCount,
			IsLiked:     comment.IsLiked,
			IsFollowing: comment.IsFollowing,
		}
		result = append(result, commentResponse)
	}
	return result
}

func NewCommentDetailResponse(comments []model.Comment) []CommentDetailResponse {
	result := []CommentDetailResponse{}
	for _, comment := range comments {

		commentDetailResponse := CommentDetailResponse{
			ID:          comment.ID,
			Content:     comment.Content,
			User:        comment.User,
			CreatedAt:   comment.CreatedAt,
			UpdatedAt:   comment.UpdatedAt,
			Replies:     NewReplyResponse(comment.Replies),
			LikeCount:   comment.LikeCount,
			ReplyCount:  comment.ReplyCount,
			IsLiked:     comment.IsLiked,
			IsFollowing: comment.IsFollowing,
		}
		result = append(result, commentDetailResponse)
	}
	return result
}
