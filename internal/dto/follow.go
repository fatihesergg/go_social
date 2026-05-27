package dto

import "github.com/google/uuid"

type SendFollowRequest struct {
	RequesterID uuid.UUID
	FollowID    uuid.UUID
}

type RespondFollowRequest struct {
	ResponderID uuid.UUID
	SenderID    uuid.UUID
}

type UnfollowRequest struct {
	RequesterID uuid.UUID
	FollowID    uuid.UUID
}

type CancelFollowRequest struct {
	RequesterID uuid.UUID
	FollowID    uuid.UUID
}
