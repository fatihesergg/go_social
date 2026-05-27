package mock

import (
	"context"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (us *MockUserStore) CreateUser(ctx context.Context, user *model.User) error {
	args := us.Called(ctx, user)
	return args.Error(0)
}

func (us *MockUserStore) UpdateUser(ctx context.Context, user *model.User) error {
	args := us.Called(ctx, user)
	return args.Error(0)
}

func (us *MockUserStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := us.Called(ctx, id)
	return args.Error(0)
}

func (us *MockUserStore) GetUsersByUsername(ctx context.Context, userName string) ([]model.User, error) {
	args := us.Called(ctx, userName)
	return args.Get(0).([]model.User), args.Error(1)
}

func (us *MockUserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := us.Called(ctx, id)
	user, ok := args.Get(0).(*model.User)
	if !ok {
		return nil, args.Error(1)
	}

	return user, args.Error(1)
}
func (us *MockUserStore) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := us.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (us *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := us.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}
