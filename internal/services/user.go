package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatihesergg/go_social/internal/appError"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/util"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type BaseUserService interface {
	Register(ctx context.Context, dto dto.CreateUserDTO) error
	Login(ctx context.Context, dto dto.LoginUserDTO) (string, error)
	ResetPassword(ctx context.Context, meID uuid.UUID, dto dto.ResetUserPasswordDTO) error
	GetUsersByUsername(ctx context.Context, username string) ([]model.User, error)
	GetUserByID(ctx context.Context, rawID string) (*model.User, error)
	GetUsersPosts(ctx context.Context, rawID string, meID uuid.UUID, limit, offset, search string) ([]model.Post, error)
}

type UserService struct {
	userStore   database.BaseUserStore
	followStore database.BaseFollowStore
	postStore   database.BasePostStore
}

func NewUserService(userstore database.BaseUserStore, followStore database.BaseFollowStore, postStore database.BasePostStore) BaseUserService {
	return &UserService{
		userStore:   userstore,
		followStore: followStore,
		postStore:   postStore,
	}
}
func (us *UserService) Register(ctx context.Context, dto dto.CreateUserDTO) error {
	user := &model.User{
		ID:       uuid.New(),
		Name:     dto.Name,
		LastName: dto.LastName,
		Email:    dto.Email,
		Avatar:   dto.Avatar,
		Username: dto.Username,
		Password: dto.Password,
	}

	existEmail, err := us.userStore.GetUserByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return appError.InternalServerError.Wrap(err)
	}
	if existEmail != nil {
		return appError.EmailExistError
	}

	existUsername, err := us.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return appError.InternalServerError.Wrap(err)
	}
	if existUsername != nil {
		return appError.UsernameExistError.Wrap(fmt.Errorf("request user username already exist"))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return appError.InternalServerError.Wrap(err)
	}
	user.Password = string(hashedPass)

	err = us.userStore.CreateUser(ctx, user)
	if err != nil {
		return appError.InternalServerError.Wrap(err)
	}

	return nil
}

func (us *UserService) GetUserByID(ctx context.Context, rawID string) (*model.User, error) {
	userID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, appError.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing userid: %w", err))
	}

	user, err := us.userStore.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.UserNotFoundError.Wrap(fmt.Errorf("user not found"))
		}
		return nil, appError.InternalServerError.Wrap(err)
	}

	return user, nil

}

func (us *UserService) Login(ctx context.Context, dto dto.LoginUserDTO) (string, error) {
	user, err := us.userStore.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", appError.UserNotFoundError.Wrap(err)
		}
		return "", appError.InternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", appError.InvalidCredentialsError.Wrap(err)
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		return "", appError.InternalServerError.Wrap(fmt.Errorf("error while creating jwt token: %w", err))
	}
	return token, nil
}

func (us *UserService) GetUsersPosts(ctx context.Context, rawID string, meID uuid.UUID, limit, offset string, query string) ([]model.Post, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		return nil, appError.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing userid: %w", err))
	}

	user, err := us.userStore.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.UserNotFoundError.Wrap(fmt.Errorf("user not found"))
		}
		return nil, appError.InternalServerError.Wrap(err)
	}

	isFollowing, err := us.followStore.IsFollowing(ctx, model.Follow{UserID: meID, FollowID: userID})

	if err != nil {
		return nil, appError.InternalServerError.Wrap(err)
	}
	if !isFollowing {
		return nil, appError.NotFollowingError
	}

	followers, err := us.followStore.GetFollowerByUserID(ctx, user.ID)
	if err != nil {
		return nil, appError.InternalServerError
	}

	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)
	//NOTE: don't loop array make query
	for i := range followers {
		followerID := followers[i].UserID
		if followerID == meID {
			posts, err := us.postStore.GetPostsByUserID(ctx, user.ID, pagination, search)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, appError.NoPostsFoundError.Wrap(fmt.Errorf("no posts found"))
				}
				return nil, appError.InternalServerError.Wrap(err)
			}

			return posts, nil
		}
	}
	return nil, appError.NoFollowersFoundError
}
func (us *UserService) ResetPassword(ctx context.Context, meID uuid.UUID, dto dto.ResetUserPasswordDTO) error {

	user, err := us.userStore.GetUserByID(ctx, meID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appError.UserNotFoundError.Wrap(fmt.Errorf("user not found"))
		}
		return appError.InternalServerError.Wrap(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.OldPassword)); err != nil {
		return appError.InvalidCredentialsError.Wrap(err)
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return appError.InternalServerError.Wrap(fmt.Errorf("error while hashing password: %w", err))
	}
	user.Password = string(hashedPass)

	err = us.userStore.UpdateUser(ctx, user)
	if err != nil {
		return appError.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) GetUsersByUsername(ctx context.Context, username string) ([]model.User, error) {

	users, err := us.userStore.GetUsersByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appError.NoUsersFoundError.Wrap(fmt.Errorf("no users found"))
		}
		return nil, appError.InternalServerError.Wrap(err)
	}

	return users, nil
}
