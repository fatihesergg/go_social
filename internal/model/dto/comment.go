package dto

import "github.com/google/uuid"

type CreateCommentDTO struct {
	PostID  uuid.UUID `json:"post_id" binding:"required,uuid"`
	Content string    `json:"content" binding:"required,alphanum,lte=200"`
	Image   string    `json:"image"`
}

type UpdateCommentDTO struct {
	Content string `json:"content" binding:"required,alphanum,lte=200"`
	Image   string `json:"image"`
}
