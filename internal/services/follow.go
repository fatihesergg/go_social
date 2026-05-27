package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatihesergg/go_social/internal/appError"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseFollowService interface {
	SendFollowRequest(ctx context.Context, dto dto.SendFollowRequest) error
	CancelFollowRequest(ctx context.Context, dto dto.CancelFollowRequest) error
	AcceptFollowRequest(ctx context.Context, dto dto.RespondFollowRequest) error
	RejectFollowRequest(ctx context.Context, dto dto.RespondFollowRequest) error
	GetFollowerByUserID(ctx context.Context, followModel model.Follow) ([]model.Follow, error)
	GetFollowingByUserID(ctx context.Context, followModel model.Follow) ([]model.Follow, error)
	UnFollowUser(ctx context.Context, dto dto.UnfollowRequest) error
}

type FollowService struct {
	followStore database.BaseFollowStore
	userStore   database.BaseUserStore
}

func NewFollowService(followStore database.BaseFollowStore, userStore database.BaseUserStore) BaseFollowService {
	return &FollowService{
		followStore: followStore,
		userStore:   userStore,
	}
}

func (fs *FollowService) SendFollowRequest(ctx context.Context, dto dto.SendFollowRequest) error {

	if dto.RequesterID == dto.FollowID {
		return appError.SelfFollowError
	}

	status, err := fs.followStore.UpsertFollowRequest(ctx, model.Follow{ID: uuid.New(), UserID: dto.RequesterID, FollowID: dto.FollowID, Status: model.Pending})
	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while getting upsert follow: %w", err))
	}

	switch status {
	case string(model.Pending):
		return appError.AlreadySendFollowRequestError
	case string(model.Accepted):
		return appError.AlreadyFollowingError
	default:
		return nil
	}
}

func (fs *FollowService) CancelFollowRequest(ctx context.Context, dto dto.CancelFollowRequest) error {

	if dto.RequesterID == dto.FollowID {
		return appError.SelfFollowError
	}

	oldStatus, newStatus, err := fs.followStore.DeleteFollowRequest(ctx, model.Follow{UserID: dto.RequesterID, FollowID: dto.FollowID}, model.Pending)

	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while cancel follow request: %w", err))
	}

	switch {
	case oldStatus == "" && newStatus == "":
		return appError.NotSendFollowRequest
	case oldStatus == string(model.Rejected) && newStatus == "":
		return appError.AlreadyRejectedFollowRequestError
	case oldStatus == string(model.Accepted) && newStatus == "":
		return appError.AlreadyAcceptedFollowRequestError
	case oldStatus == string(model.Pending) && newStatus == string(model.Pending):
		return nil
	}

	return appError.InternalServerError.Wrap(fmt.Errorf("unknown value for oldStatus %s and newStatus %s", oldStatus, newStatus))

}

func (fs *FollowService) AcceptFollowRequest(ctx context.Context, dto dto.RespondFollowRequest) error {

	if dto.ResponderID == dto.SenderID {
		return appError.SelfFollowError
	}

	oldStatus, newStatus, err := fs.followStore.UpdateFollowStatus(ctx, model.Follow{UserID: dto.SenderID, FollowID: dto.ResponderID, Status: model.Accepted})
	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while getting follow status: %w", err))
	}

	switch {
	case oldStatus == "" && newStatus == "":
		return appError.NotSendFollowRequest
	case oldStatus == string(model.Accepted) && newStatus == "":
		return appError.AlreadyAcceptedFollowRequestError
	case oldStatus == string(model.Rejected) && newStatus == "":
		return appError.AlreadyRejectedFollowRequestError
	case oldStatus == string(model.Pending) && newStatus == string(model.Accepted):
		return nil
	}

	return appError.InternalServerError.Wrap(fmt.Errorf("unknown values"))

}
func (fs *FollowService) RejectFollowRequest(ctx context.Context, dto dto.RespondFollowRequest) error {

	if dto.ResponderID == dto.SenderID {
		return appError.SelfFollowError
	}

	oldStatus, newStatus, err := fs.followStore.UpdateFollowStatus(ctx, model.Follow{UserID: dto.SenderID, FollowID: dto.ResponderID, Status: model.Rejected})
	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while getting follow status: %w", err))
	}

	switch {
	case oldStatus == "" && newStatus == "":
		return appError.NotSendFollowRequest
	case oldStatus == string(model.Accepted) && newStatus == "":
		return appError.AlreadyAcceptedFollowRequestError
	case oldStatus == string(model.Rejected) && newStatus == "":
		return appError.AlreadyRejectedFollowRequestError
	case oldStatus == string(model.Pending) && newStatus == string(model.Rejected):
		return nil
	}
	return appError.InternalServerError.Wrap(fmt.Errorf("unknown value for oldStatus %s and newStatus %s", oldStatus, newStatus))

}
func (fs *FollowService) GetFollowerByUserID(ctx context.Context, followModel model.Follow) ([]model.Follow, error) {

	isFollowing, err := fs.followStore.IsFollowing(ctx, followModel)
	if err != nil {
		return nil, appError.InternalServerError.Wrap(fmt.Errorf("error while checking is following: %w", err))
	}

	if !isFollowing {
		return nil, appError.NotFollowingError
	}

	followers, err := fs.followStore.GetFollowerByUserID(ctx, followModel.FollowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.NoFollowersFoundError.Wrap(err)
		}
		return nil, appError.InternalServerError.Wrap(err)
	}

	return followers, nil
}
func (fs *FollowService) GetFollowingByUserID(ctx context.Context, followModel model.Follow) ([]model.Follow, error) {

	isFollowing, err := fs.followStore.IsFollowing(ctx, followModel)
	if err != nil {
		return nil, appError.InternalServerError.Wrap(fmt.Errorf("error while checking is following: %w", err))
	}

	if !isFollowing {
		return nil, appError.NotFollowingError
	}

	followings, err := fs.followStore.GetFollowingByUserID(ctx, followModel.FollowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.NoFollowingsFoundError.Wrap(fmt.Errorf("no followings found"))
		}
		return nil, appError.InternalServerError.Wrap(err)
	}

	return followings, nil
}

func (fs *FollowService) UnFollowUser(ctx context.Context, dto dto.UnfollowRequest) error {

	oldStatus, newStatus, err := fs.followStore.DeleteFollowRequest(ctx, model.Follow{UserID: dto.RequesterID, FollowID: dto.FollowID}, model.Accepted)

	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while unfollow user: %w", err))
	}

	switch {
	case oldStatus == "" && newStatus == "" || oldStatus == string(model.Rejected) && newStatus == "":
		return appError.NotFollowingError
	case oldStatus == string(model.Pending) && newStatus == "":
		return appError.FollowNotAccepted
	case oldStatus == string(model.Accepted) && newStatus == string(model.Accepted):
		return nil

	}

	return appError.InternalServerError.Wrap(fmt.Errorf("unknown value for oldStatus %s and newStatus %s", oldStatus, newStatus))
}
