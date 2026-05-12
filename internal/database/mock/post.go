package mock

import (
	"context"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostStore struct {
	mock.Mock
}

func (ps *MockPostStore) HasAccessToPost(ctx context.Context, userID, postID uuid.UUID) (bool, error) {
	args := ps.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}
func (ps *MockPostStore) GetPostByID(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	args := ps.Called(ctx, postID)
	model, ok := args.Get(0).(*model.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return model, args.Error(1)
}
func (ps *MockPostStore) GetPosts(ctx context.Context, pagination database.Pagination, search database.Search, userID uuid.UUID) ([]model.Post, error) {
	args := ps.Called(ctx, pagination, search, userID)
	arr, ok := args.Get(0).([]model.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return arr, args.Error(1)
}
func (ps *MockPostStore) GetPostsByTag(ctx context.Context, pagination database.Pagination, tag string, userID uuid.UUID) ([]model.Post, error) {
	args := ps.Called(ctx, pagination, tag, userID)
	arr, ok := args.Get(0).([]model.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return arr, args.Error(1)
}
func (ps *MockPostStore) GetPostDetailsByID(ctx context.Context, postID, userID uuid.UUID) (*model.Post, error) {
	args := ps.Called(ctx, postID, userID)
	model, ok := args.Get(0).(*model.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return model, args.Error(1)
}
func (ps *MockPostStore) GetPostsByUserID(ctx context.Context, userID uuid.UUID, pagination database.Pagination, search database.Search) ([]model.Post, error) {
	args := ps.Called(ctx, userID, pagination, search)
	arr, ok := args.Get(0).([]model.Post)
	if !ok {
		return nil, args.Error(1)
	}
	return arr, args.Error(1)
}
func (ps *MockPostStore) CreatePost(ctx context.Context, post *model.Post) error {
	args := ps.Called(ctx, post)
	return args.Error(1)

}
func (ps *MockPostStore) UpdatePost(ctx context.Context, post *model.Post) error {
	args := ps.Called(ctx, post)
	return args.Error(1)

}
func (ps *MockPostStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	args := ps.Called(ctx, id)
	return args.Error(1)
}
