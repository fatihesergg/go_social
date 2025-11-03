package dto

import (
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type CreateReply struct {
	Message string `json:"message",binding:"required,lte=100"`
}

type UpdateReply struct {
	Message string `json:"message",binding:"required,lte=100"`
}

type ReplyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Message     string     `json:"message"`
	User        model.User `json:"user"`
	ReplyCount  int        `json:"total_reply"`
	LikeCount   int        `json:"total_like"`
	IsLiked     bool       `json:"is_liked"`
	IsFollowing bool       `json:"is_followed"`
}

func NewReplyResponse(replies []model.Reply) []ReplyResponse {
	result := []ReplyResponse{}
	for _, reply := range replies {
		replyResponse := ReplyResponse{
			ID:          reply.ID,
			Message:     reply.Message,
			User:        reply.User,
			ReplyCount:  reply.ReplyCount,
			LikeCount:   reply.LikeCount,
			IsFollowing: reply.IsFollowing,
			IsLiked:     reply.IsLiked,
		}
		result = append(result, replyResponse)
	}
	return result
}
