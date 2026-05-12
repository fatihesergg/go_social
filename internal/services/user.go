package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
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
	GetFollowerByUserID(ctx context.Context, rawID string) ([]model.Follow, error)
	GetFollowingByUserID(ctx context.Context, rawID string) ([]model.Follow, error)
	FollowUser(ctx context.Context, userID uuid.UUID, followID string) error
	UnFollowUser(ctx context.Context, userID uuid.UUID, followID string) error
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
		return errors.InternalServerError.Wrap(err)
	}
	if existEmail != nil {
		return errors.EmailExistError.Wrap(fmt.Errorf("Request user email already exist"))
	}

	existUsername, err := us.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil && err != sql.ErrNoRows {
		return errors.InternalServerError.Wrap(err)
	}
	if existUsername != nil {
		return errors.UsernameExistError.Wrap(fmt.Errorf("Request user username already exist"))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Invalid credentials: %w", err)).Wrap(err)
	}
	user.Password = string(hashedPass)

	err = us.userStore.CreateUser(ctx, user)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil
}

func (us *UserService) GetUserByID(ctx context.Context, rawID string) (*model.User, error) {
	userID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	user, err := us.userStore.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return user, nil

}

func (us *UserService) Login(ctx context.Context, dto dto.LoginUserDTO) (string, error) {
	user, err := us.userStore.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return "", errors.InternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", errors.InvalidCredentialsError.Wrap(fmt.Errorf("Invalid credentials: %w", err))
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		return "", errors.InternalServerError.Wrap(err).Wrap(fmt.Errorf("Error while creating jwt token: %w", err))
	}
	return token, nil
}
func (us *UserService) GetFollowerByUserID(ctx context.Context, rawID string) ([]model.Follow, error) {

	userID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	followers, err := us.followStore.GetFollowerByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoFollowersFoundError.Wrap(fmt.Errorf("No followers found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return followers, nil
}
func (us *UserService) GetFollowingByUserID(ctx context.Context, rawID string) ([]model.Follow, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	followings, err := us.followStore.GetFollowingByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoFollowingsFoundError.Wrap(fmt.Errorf("No followings found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return followings, nil
}
func (us *UserService) FollowUser(ctx context.Context, userID uuid.UUID, followID string) error {

	followUserID, err := uuid.Parse(followID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing followUserID: %w", err))
	}

	followings, err := us.followStore.GetFollowingByUserID(ctx, userID)

	if err != nil && err != sql.ErrNoRows {
		return errors.InternalServerError.Wrap(err)
	}

	isFollowing := false
	for _, follow := range followings {
		if follow.FollowID == followUserID {
			isFollowing = true
		}
	}

	if isFollowing {
		return errors.AlreadyFollowingError.Wrap(fmt.Errorf("Request user already follow this user"))
	}
	follow := model.Follow{
		ID:       uuid.New(),
		UserID:   userID,
		FollowID: followUserID,
	}

	err = us.followStore.FollowUser(ctx, follow)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) UnFollowUser(ctx context.Context, userID uuid.UUID, followID string) error {

	unfUser, err := uuid.Parse(followID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing unfUser: %w", err))
	}

	isFollowing, err := us.followStore.IsFollowing(ctx, model.Follow{UserID: userID, FollowID: unfUser})

	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while checking user following: %w", err))
	}

	if !isFollowing {
		return errors.NotFollowingError
	}
	follow := model.Follow{
		UserID:   userID,
		FollowID: unfUser,
	}

	err = us.followStore.UnFollowUser(ctx, follow)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) GetUsersPosts(ctx context.Context, rawID string, meID uuid.UUID, limit, offset string, query string) ([]model.Post, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	user, err := us.userStore.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	isFollowing, err := us.followStore.IsFollowing(ctx, model.Follow{UserID: meID, FollowID: userID})

	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}
	if !isFollowing {
		return nil, errors.NotFollowingError
	}

	followers, err := us.followStore.GetFollowerByUserID(ctx, user.ID)
	if err != nil {
		return nil, errors.InternalServerError
	}

	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)
	//NOTE: don't loop array make query
	for i := range followers {
		followerID := followers[i].UserID
		if followerID == meID {
			posts, err := us.postStore.GetPostsByUserID(ctx, user.ID, pagination, search)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
				}
				return nil, errors.InternalServerError.Wrap(err)
			}

			return posts, nil
		}
	}
	return nil, errors.NoFollowersFoundError
}
func (us *UserService) ResetPassword(ctx context.Context, meID uuid.UUID, dto dto.ResetUserPasswordDTO) error {

	user, err := us.userStore.GetUserByID(ctx, meID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.OldPassword)); err != nil {
		return errors.InvalidCredentialsError.Wrap(fmt.Errorf("Invalid credentials"))
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while hashing password: %w", err)).Wrap(err)
	}
	user.Password = string(hashedPass)

	err = us.userStore.UpdateUser(ctx, user)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) GetUsersByUsername(ctx context.Context, username string) ([]model.User, error) {

	users, err := us.userStore.GetUsersByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoUsersFoundError.Wrap(fmt.Errorf("No users found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return users, nil
}
