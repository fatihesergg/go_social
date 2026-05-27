package model

import "github.com/google/uuid"

type Follow struct {
	ID       uuid.UUID       `json:"-"`
	UserID   uuid.UUID       `json:"user_id"`
	FollowID uuid.UUID       `json:"follow_id"`
	Status   FollowReqStatus `json:"-"`
}

type FollowReqStatus string

const (
	Accepted FollowReqStatus = "accepted"
	Pending  FollowReqStatus = "pending"
	Rejected FollowReqStatus = "rejected"
)
