package services

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/util"
	"go.uber.org/zap"

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
	logger  *zap.Logger
}

func NewUserService(storage *database.Storage, logger *zap.Logger) BaseUserService {
	return &UserService{storage: storage, logger: logger.Named("user_service")}
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
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if existEmail != nil {
		us.logger.Error("Request user email already exist")
		return errors.EmailExistError
	}

	existUsername, err := us.storage.UserStore.GetUserByUsername(user.Username)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if existUsername != nil {
		us.logger.Error("Request user username already exist")
		return errors.UsernameExistError
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		us.logger.Error("Invalid credentials")
		return errors.InternalServerError.Wrap(err)
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
		logger.Error("Error while parsing userID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	user, err := us.storage.UserStore.GetUserByID(userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if user == nil {
		return nil, errors.UserNotFoundError
	}
	return user, nil

}

func (us *UserService) Login(dto dto.LoginUserDTO) (string, error) {
	user, err := us.storage.UserStore.GetUserByEmail(dto.Email)
	if err != nil {
		return "", errors.InternalServerError
	}
	if user == nil {
		return "", errors.UserNotFoundError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		us.logger.Error("Invalid credentials")
		return "", errors.InvalidCredentialsError
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		us.logger.Error("Error while creating jwt token", zap.Error(err))
		return "", errors.InternalServerError.Wrap(err)
	}
	return token, nil
}
func (us *UserService) GetFollowerByUserID(rawID string) ([]model.Follow, error) {

	userID, err := uuid.Parse(rawID)
	if err != nil {
		logger.Error("Error while parsing userID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	followers, err := us.storage.FollowStore.GetFollowerByUserID(userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if len(followers) == 0 {
		return nil, errors.NoFollowersFoundError
	}
	return followers, nil
}
func (us *UserService) GetFollowingByUserID(rawID string) ([]model.Follow, error) {
	userID, err := uuid.Parse(rawID)

	if err != nil {
		us.logger.Error("Error while parsing userID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	followings, err := us.storage.FollowStore.GetFollowingByUserID(userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if len(followings) == 0 {
		return nil, errors.NoFollowingsFoundError
	}
	return followings, nil
}
func (us *UserService) FollowUser(userID uuid.UUID, followID string) error {

	followUserID, err := uuid.Parse(followID)
	if err != nil {
		logger.Error("Error while parsing followUserID", zap.Error(err))
		return errors.InvalidIDFormatError
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
		us.logger.Error("Request user already follow this user")
		return errors.AlreadyFollowingError
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
		logger.Error("Error while parsing unfUser", zap.Error(err))
		return errors.InvalidIDFormatError
	}

	followings, err := us.storage.FollowStore.GetFollowingByUserID(userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	isFollowing := false
	for _, follow := range followings {
		if follow.FollowID == unfUser {
			isFollowing = true
		}
	}
	if !isFollowing {
		us.logger.Error("User has not following this user yet")
		return errors.NotFollowingError
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
		us.logger.Error("Error while parsing userID")
		return nil, errors.InvalidIDFormatError
	}

	user, err := us.storage.UserStore.GetUserByID(userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if user == nil {
		return nil, errors.UserNotFoundError
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
				return nil, errors.InternalServerError.Wrap(err)
			}
			if posts == nil {
				return nil, errors.NoPostsFoundError
			}
			result := dto.NewAllPostResponse(posts)
			return result, nil
		}
	}
	us.logger.Error("Request user has not following this user yet")
	return nil, errors.NotFollowingError
}
func (us *UserService) ResetPassword(meID uuid.UUID, dto dto.ResetUserPasswordDTO) error {

	user, err := us.storage.UserStore.GetUserByID(meID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if user == nil {
		return errors.UserNotFoundError
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.OldPassword)); err != nil {
		us.logger.Error("Invalid credentials")
		return errors.InvalidCredentialsError
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		us.logger.Error("Error while hashing password", zap.Error(err))
		return errors.InternalServerError.Wrap(err)
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
		return nil, errors.InternalServerError.Wrap(err)
	}

	if users == nil {
		return nil, errors.UserNotFoundError
	}
	return users, nil
}
