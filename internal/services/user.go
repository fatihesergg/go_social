package services

import (
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
	Register(dto dto.CreateUserDTO) error
	Login(dto dto.LoginUserDTO) (string, error)
	ResetPassword(meID uuid.UUID, dto dto.ResetUserPasswordDTO) error
	GetUsersByUsername(username string) ([]model.User, error)
	GetUserByID(rawID string) (*model.User, error)
	GetFollowerByUserID(rawID string) ([]model.Follow, error)
	GetFollowingByUserID(rawID string) ([]model.Follow, error)
	FollowUser(userID uuid.UUID, followID string) error
	UnFollowUser(userID uuid.UUID, followID string) error
	GetUsersPosts(rawID string, meID uuid.UUID, limit, offset, search string) ([]dto.AllPostResponse, error)
}

type UserService struct {
	storage *database.Storage
}

func NewUserService(storage *database.Storage) BaseUserService {
	return &UserService{storage: storage}
}
func (us *UserService) Register(dto dto.CreateUserDTO) error {
	user := &model.User{
		ID:       uuid.New(),
		Name:     dto.Name,
		LastName: dto.LastName,
		Email:    dto.Email,
		Avatar:   dto.Avatar,
		Username: dto.Username,
		Password: dto.Password,
	}

	existEmail, err := us.storage.UserStore.GetUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		return errors.InternalServerError.Wrap(err)
	}
	if existEmail != nil {
		return errors.EmailExistError.Wrap(fmt.Errorf("Request user email already exist"))
	}

	existUsername, err := us.storage.UserStore.GetUserByUsername(user.Username)
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

	err = us.storage.UserStore.CreateUser(user)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil
}

func (us *UserService) GetUserByID(rawID string) (*model.User, error) {
	userID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	user, err := us.storage.UserStore.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return user, nil

}

func (us *UserService) Login(dto dto.LoginUserDTO) (string, error) {
	user, err := us.storage.UserStore.GetUserByEmail(dto.Email)
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
func (us *UserService) GetFollowerByUserID(rawID string) ([]model.Follow, error) {

	userID, err := uuid.Parse(rawID)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	followers, err := us.storage.FollowStore.GetFollowerByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoFollowersFoundError.Wrap(fmt.Errorf("No followers found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return followers, nil
}
func (us *UserService) GetFollowingByUserID(rawID string) ([]model.Follow, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	followings, err := us.storage.FollowStore.GetFollowingByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoFollowingsFoundError.Wrap(fmt.Errorf("No followings found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return followings, nil
}
func (us *UserService) FollowUser(userID uuid.UUID, followID string) error {

	followUserID, err := uuid.Parse(followID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing followUserID: %w", err))
	}

	followings, err := us.storage.FollowStore.GetFollowingByUserID(userID)

	if err != nil {
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

	err = us.storage.FollowStore.FollowUser(follow)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) UnFollowUser(userID uuid.UUID, followID string) error {

	unfUser, err := uuid.Parse(followID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing unfUser: %w", err))
	}

	followings, err := us.storage.FollowStore.GetFollowingByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NoFollowingsFoundError.Wrap(fmt.Errorf("No followings found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	isFollowing := false
	for _, follow := range followings {
		if follow.FollowID == unfUser {
			isFollowing = true
		}
	}
	if !isFollowing {
		return errors.NotFollowingError.Wrap(fmt.Errorf("User has not following this user yet: %w", err))
	}
	follow := model.Follow{
		UserID:   userID,
		FollowID: unfUser,
	}

	err = us.storage.FollowStore.UnFollowUser(follow)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) GetUsersPosts(rawID string, meID uuid.UUID, limit, offset string, query string) ([]dto.AllPostResponse, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing userID: %w", err))
	}

	user, err := us.storage.UserStore.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.UserNotFoundError.Wrap(fmt.Errorf("User not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	followers, err := us.storage.FollowStore.GetFollowerByUserID(user.ID)
	if err != nil {
		return nil, errors.InternalServerError
	}

	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)
	for i := range followers {
		followerID := followers[i].ID
		if followerID == userID {
			posts, err := us.storage.PostStore.GetPostsByUserID(user.ID, pagination, search)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
				}
				return nil, errors.InternalServerError.Wrap(err)
			}

			result := dto.NewAllPostResponse(posts)
			return result, nil
		}
	}
	return nil, errors.NotFollowingError.Wrap(fmt.Errorf("Request user has not following this user yet: %w", err))
}
func (us *UserService) ResetPassword(meID uuid.UUID, dto dto.ResetUserPasswordDTO) error {

	user, err := us.storage.UserStore.GetUserByID(meID)
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

	err = us.storage.UserStore.UpdateUser(user)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (us *UserService) GetUsersByUsername(username string) ([]model.User, error) {

	users, err := us.storage.UserStore.GetUsersByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoUsersFoundError.Wrap(fmt.Errorf("No users found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return users, nil
}
