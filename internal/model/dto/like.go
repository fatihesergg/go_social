package dto

import "github.com/google/uuid"

type CreatePostLikeDTO struct {
	PostID uuid.UUID `json:"post_id" binding:"required,uuid"`
}

type CreateCommentLikeDTO struct {
	CommentID uuid.UUID `json:"post_id" binding:"required,uuid"`
}
