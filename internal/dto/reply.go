package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type CreateReply struct {
	Message string `json:"message",binding:"required,alphanum,lte=100"`
}

type UpdateReply struct {
	Message string `json:"message",binding:"required,alphanum,lte=100"`
}

type ReplyResponse struct {
	ID      uuid.UUID  `json:"id"`
	Message string     `json:"message"`
	User    model.User `json:"user"`
}

func NewReplyResponse(replies []model.Reply) []ReplyResponse {
	result := []ReplyResponse{}
	for _, reply := range replies {
		replyResponse := ReplyResponse{
			ID:      reply.ID,
			Message: reply.Message,
			User:    reply.User,
		}
		result = append(result, replyResponse)
	}
	return result
}
