package dto

import "github.com/google/uuid"

type PostLikedEvent struct {
	LikerID uuid.UUID `json:"liker_id"`
	PostID  uuid.UUID `json:"post_id"`
}
