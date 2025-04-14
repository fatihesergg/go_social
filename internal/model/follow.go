package model

type Follow struct {
	ID       int64 `json:"-"`
	UserID   int64 `json:"user_id"`
	FollowID int64 `json:"follow_id"`
}
