package mock

import (
	"context"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFollowStore struct {
	mock.Mock
}

func (fs *MockFollowStore) GetFollowerByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	args := fs.Called(ctx, userID)
	follows, ok := args.Get(0).([]model.Follow)
	if !ok {
		return nil, args.Error(1)
	}
	return follows, args.Error(1)
}

func (fs *MockFollowStore) GetFollowingByUserID(ctx context.Context, userID uuid.UUID) ([]model.Follow, error) {
	args := fs.Called(ctx, userID)
	follows, ok := args.Get(0).([]model.Follow)
	if !ok {
		return nil, args.Error(1)
	}
	return follows, args.Error(1)
}
func (fs *MockFollowStore) UpsertFollowRequest(ctx context.Context, followModel model.Follow) (string, error) {
	args := fs.Called(ctx, followModel)
	return args.String(0), args.Error(1)
}

func (fs *MockFollowStore) DeleteFollowRequest(ctx context.Context, model model.Follow, status model.FollowReqStatus) (string, string, error) {
	args := fs.Called(ctx, model, status)
	return args.String(0), args.String(1), args.Error(1)
}

func (fs *MockFollowStore) UpdateFollowStatus(ctx context.Context, model model.Follow) (string, string, error) {
	args := fs.Called(ctx, model)
	return args.String(0), args.String(1), args.Error(1)
}
func (fs *MockFollowStore) IsFollowing(ctx context.Context, model model.Follow) (bool, error) {
	args := fs.Called(ctx, model)
	return args.Bool(0), args.Error(1)
}
