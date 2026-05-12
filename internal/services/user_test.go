package services

import (
	"database/sql"
	"fmt"
	"testing"

	db_mock "github.com/fatihesergg/go_social/internal/database/mock"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setup(t *testing.T) (*db_mock.MockUserStore, *db_mock.MockFollowStore, *db_mock.MockPostStore, BaseUserService) {
	t.Helper()
	userStore := new(db_mock.MockUserStore)
	followStore := new(db_mock.MockFollowStore)
	postStore := new(db_mock.MockPostStore)
	userService := NewUserService(userStore, followStore, postStore)
	return userStore, followStore, postStore, userService
}

func TestRegister_DBError(t *testing.T) {

	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	err := service.Register(t.Context(), dto.CreateUserDTO{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())

}

func TestRegister_EmailExist(t *testing.T) {
	mockStore, _, _, service := setup(t)

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{}, nil)
	err := service.Register(t.Context(), dto.CreateUserDTO{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.EmailExistError.Error())
}

func TestRegister_UsernameExist(t *testing.T) {
	mockStore, _, _, service := setup(t)

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("GetUserByUsername", mock.Anything, mock.Anything).Return(&model.User{}, nil)
	err := service.Register(t.Context(), dto.CreateUserDTO{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.UsernameExistError.Error())

}
func TestRegister_GetUserByUsernameFail(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	err := service.Register(t.Context(), dto.CreateUserDTO{})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}
func TestRegister_CreateUserFail(t *testing.T) {
	mockStore, _, _, service := setup(t)

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("CreateUser", mock.Anything, mock.Anything).Return(fmt.Errorf("Unexpected Error"))

	err := service.Register(t.Context(), dto.CreateUserDTO{})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestRegister_CreateUserSuccess(t *testing.T) {
	mockStore, _, _, service := setup(t)

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	err := service.Register(t.Context(), dto.CreateUserDTO{})

	assert.NoError(t, err)

}

func TestGetUserByID_InvalidUUID(t *testing.T) {
	mockStore, _, _, service := setup(t)

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("GetUserByUsername", mock.Anything, mock.Anything).Return(nil, nil)
	mockStore.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	user, err := service.GetUserByID(t.Context(), "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())
	assert.Nil(t, user)
}
func TestGetUserByID_GetUserByIDNoRows(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	user, err := service.GetUserByID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, err.Error(), errors.UserNotFoundError.Error())
}

func TestGetUserByID_GetUserByIDFail(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	user, err := service.GetUserByID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}
func TestGetUserByID_GetUserByIDSuccess(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Username: "test_user"}, nil)

	user, err := service.GetUserByID(t.Context(), uuid.NewString())

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Username, "test_user")

}

func TestLogin_GetUserByEmailNoRows(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	str, err := service.Login(t.Context(), dto.LoginUserDTO{})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.UserNotFoundError.Error())
	assert.Empty(t, str)
}

func TestLogin_GetUserByEmailFail(t *testing.T) {
	mockStore, _, _, service := setup(t)
	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	str, err := service.Login(t.Context(), dto.LoginUserDTO{})
	assert.Error(t, err)
	assert.Equal(t, err, errors.InternalServerError)
	assert.Empty(t, str)
}
func TestLogin_HashFail(t *testing.T) {
	mockStore, _, _, service := setup(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

	if err != nil {
		t.Fail()
	}

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{Password: string(hash)}, nil)

	str, err := service.Login(t.Context(), dto.LoginUserDTO{Password: "12"})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidCredentialsError.Error())
	assert.Empty(t, str)
}

func TestLogin_Success(t *testing.T) {
	mockStore, _, _, service := setup(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)

	if err != nil {
		t.Fail()
	}
	t.Setenv("JWT_SECRET", "123")

	util.ApiConfig = &util.Config{JWTSecret: "123"}

	mockStore.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&model.User{ID: uuid.New(), Password: string(hash)}, nil)
	str, err := service.Login(t.Context(), dto.LoginUserDTO{Password: "123"})
	assert.NoError(t, err)
	assert.NotEmpty(t, str)

}

func TestGetFollowerByUserID_InvalidUUID(t *testing.T) {
	_, _, _, service := setup(t)

	follows, err := service.GetFollowerByUserID(t.Context(), "")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())
	assert.Nil(t, follows)
}

func TestGetFollowerByUserID_GetFollowerByUserIDNoRows(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	follows, err := service.GetFollowerByUserID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.NoFollowersFoundError.Error())
	assert.Nil(t, follows)

}
func TestGetFollowerByUserID_GetFollowerByUserFail(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("UnExpected error"))

	follows, err := service.GetFollowerByUserID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
	assert.Nil(t, follows)
}
func TestGetFollowerByUserID_Success(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return([]model.Follow{}, nil)

	follows, err := service.GetFollowerByUserID(t.Context(), uuid.NewString())

	assert.NoError(t, err)
	assert.NotNil(t, follows)
}

func TestGetFollowingByUserID_InvalidUUID(t *testing.T) {
	_, _, _, service := setup(t)

	follows, err := service.GetFollowingByUserID(t.Context(), "")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())
	assert.Nil(t, follows)
}
func TestGetFollowingByUserID_GetFollowingByUserIDNoRows(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	follows, err := service.GetFollowingByUserID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Error(t, err, errors.InvalidIDFormatError.Error())
	assert.Nil(t, follows)
}
func TestGetFollowingByUserID_GetFollowingByUserIDFail(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	follows, err := service.GetFollowingByUserID(t.Context(), uuid.NewString())

	assert.Error(t, err)
	assert.Error(t, err, errors.InternalServerError.Error())
	assert.Nil(t, follows)
}
func TestGetFollowingByUserID_Success(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return([]model.Follow{}, nil)

	follows, err := service.GetFollowingByUserID(t.Context(), uuid.NewString())

	assert.NoError(t, err)
	assert.NotNil(t, follows)
}

func TestFollowUser_InvalidUUID(t *testing.T) {
	_, _, _, service := setup(t)

	err := service.FollowUser(t.Context(), uuid.Nil, "")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())

}

func TestFollowUser_GetFollowingByUserIDFail(t *testing.T) {
	_, followStore, _, service := setup(t)

	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))
	err := service.FollowUser(t.Context(), uuid.New(), uuid.NewString())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestFollowUser_AlreadyFollowing(t *testing.T) {
	_, followStore, _, service := setup(t)
	user, follow := uuid.New(), uuid.New()

	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return([]model.Follow{model.Follow{UserID: user, FollowID: follow}}, nil)
	err := service.FollowUser(t.Context(), user, follow.String())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.AlreadyFollowingError.Error())
}

func TestFollowUser_FollowUserFail(t *testing.T) {
	_, followStore, _, service := setup(t)
	user, follow := uuid.New(), uuid.New()

	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return([]model.Follow{}, nil)
	followStore.On("FollowUser", mock.Anything, mock.Anything).Return(fmt.Errorf("unexpected error"))
	err := service.FollowUser(t.Context(), user, follow.String())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}
func TestFollowUser_Success(t *testing.T) {
	_, followStore, _, service := setup(t)
	user, follow := uuid.New(), uuid.New()

	followStore.On("GetFollowingByUserID", mock.Anything, mock.Anything).Return([]model.Follow{}, nil)
	followStore.On("FollowUser", mock.Anything, mock.Anything).Return(nil)
	err := service.FollowUser(t.Context(), user, follow.String())

	assert.NoError(t, err)
}
func TestUnFollowUser_InvalidUUID(t *testing.T) {
	_, _, _, service := setup(t)

	err := service.UnFollowUser(t.Context(), uuid.New(), "")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())
}

func TestUnFollowUser_IsFollowingError(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(false, fmt.Errorf("unexpected error"))
	err := service.UnFollowUser(t.Context(), uuid.New(), uuid.Nil.String())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestUnFollowUser_NotFollowing(t *testing.T) {
	_, followStore, _, service := setup(t)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)

	err := service.UnFollowUser(t.Context(), uuid.New(), uuid.Nil.String())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.NotFollowingError.Error())
}
func TestUnFollowUser_UnFollowUserFail(t *testing.T) {
	_, followStore, _, service := setup(t)
	user, follow := uuid.New(), uuid.New()
	followStore.On("UnFollowUser", mock.Anything, mock.Anything).Return(fmt.Errorf("unexpected error"))
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)

	err := service.UnFollowUser(t.Context(), user, follow.String())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}
func TestUnFollowUser_Success(t *testing.T) {
	_, followStore, _, service := setup(t)
	user, follow := uuid.New(), uuid.New()
	followStore.On("UnFollowUser", mock.Anything, mock.Anything).Return(nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)

	err := service.UnFollowUser(t.Context(), user, follow.String())

	assert.NoError(t, err)

}

func TestGetUsersPosts_InvalidUUID(t *testing.T) {
	_, _, _, service := setup(t)

	posts, err := service.GetUsersPosts(t.Context(), "", uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InvalidIDFormatError.Error())
}
func TestGetUsersPosts_GetUserByIDNoRows(t *testing.T) {
	userStore, _, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.UserNotFoundError.Error())
}

func TestGetUsersPosts_GetUserByIDFail(t *testing.T) {
	userStore, _, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))
	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestGetUsersPosts_GetFollowerByUserIDFail(t *testing.T) {
	userStore, followStore, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{ID: uuid.New()}, nil)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)

	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestGetUsersPosts_GetFollowerByUserIDNoRows(t *testing.T) {
	userStore, followStore, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{ID: uuid.New()}, nil)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)

	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestGetUsersPosts_NotFollowing(t *testing.T) {
	userStore, followStore, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(false, nil)
	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.NotFollowingError.Error())
}
func TestGetUsersPosts_IsFollowingFail(t *testing.T) {
	userStore, followStore, _, service := setup(t)

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(false, fmt.Errorf("unexpected error"))
	posts, err := service.GetUsersPosts(t.Context(), uuid.NewString(), uuid.New(), "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestGetUsersPosts_GetPostsByUserIDNoRows(t *testing.T) {
	userStore, followStore, postStore, service := setup(t)
	followID := uuid.New()
	userID := uuid.New()

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{ID: followID}, nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return([]model.Follow{model.Follow{UserID: userID, FollowID: followID}}, nil)
	postStore.On("GetPostsByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)
	posts, err := service.GetUsersPosts(t.Context(), followID.String(), userID, "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.NoPostsFoundError.Error())
}
func TestGetUsersPosts_GetPostsByUserIDFail(t *testing.T) {
	userStore, followStore, postStore, service := setup(t)
	followID := uuid.New()
	userID := uuid.New()

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{ID: followID}, nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return([]model.Follow{model.Follow{UserID: userID, FollowID: followID}}, nil)
	postStore.On("GetPostsByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))
	posts, err := service.GetUsersPosts(t.Context(), followID.String(), userID, "10", "0", "")

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}

func TestGetUsersPosts_GetPostsByUserSuccess(t *testing.T) {
	userStore, followStore, postStore, service := setup(t)
	followID := uuid.New()
	userID := uuid.New()

	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{ID: followID}, nil)
	followStore.On("IsFollowing", mock.Anything, mock.Anything).Return(true, nil)
	followStore.On("GetFollowerByUserID", mock.Anything, mock.Anything).Return([]model.Follow{model.Follow{UserID: userID, FollowID: followID}}, nil)
	postStore.On("GetPostsByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.Post{model.Post{}}, nil)
	posts, err := service.GetUsersPosts(t.Context(), followID.String(), userID, "10", "0", "")

	assert.NoError(t, err)
	assert.NotEmpty(t, posts)
}

func TestGetUsersByUsername_GetUsersByUsernameNoRows(t *testing.T) {
	userStore, _, _, service := setup(t)
	userStore.On("GetUsersByUsername", mock.Anything, mock.Anything).Return([]model.User{}, sql.ErrNoRows)

	users, err := service.GetUsersByUsername(t.Context(), "")

	assert.Error(t, err)
	assert.Empty(t, users)
	assert.Equal(t, err.Error(), errors.NoUsersFoundError.Error())
}

func TestGetUsersByUsername_GetUsersByUsernameFail(t *testing.T) {
	userStore, _, _, service := setup(t)
	userStore.On("GetUsersByUsername", mock.Anything, mock.Anything).Return([]model.User{}, fmt.Errorf("unexpected error"))

	users, err := service.GetUsersByUsername(t.Context(), "")

	assert.Error(t, err)
	assert.Empty(t, users)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())
}
func TestGetUsersByUsername_Success(t *testing.T) {
	userStore, _, _, service := setup(t)
	userStore.On("GetUsersByUsername", mock.Anything, mock.Anything).Return([]model.User{model.User{}}, nil)

	users, err := service.GetUsersByUsername(t.Context(), "")

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
}

func TestResetPassword_GetUserByIDNoRows(t *testing.T) {
	userStore, _, _, service := setup(t)
	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	err := service.ResetPassword(t.Context(), uuid.New(), dto.ResetUserPasswordDTO{})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.UserNotFoundError.Error())

}
func TestResetPassword_GetUserByIDFail(t *testing.T) {
	userStore, _, _, service := setup(t)
	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("unexpected error"))

	err := service.ResetPassword(t.Context(), uuid.New(), dto.ResetUserPasswordDTO{})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())

}

func TestResetPassword_InvalidCredentials(t *testing.T) {
	userStore, _, _, service := setup(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("12"), bcrypt.DefaultCost)
	if err != nil {
		t.Fail()
	}
	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Password: string(hash)}, nil)

	err = service.ResetPassword(t.Context(), uuid.New(), dto.ResetUserPasswordDTO{OldPassword: "123"})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InvalidCredentialsError.Error())

}

func TestResetPassword_UpdateUserFail(t *testing.T) {
	userStore, _, _, service := setup(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fail()
	}
	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Password: string(hash)}, nil)
	userStore.On("UpdateUser", mock.Anything, mock.Anything).Return(fmt.Errorf("unexpected error"))

	err = service.ResetPassword(t.Context(), uuid.New(), dto.ResetUserPasswordDTO{OldPassword: "123"})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), errors.InternalServerError.Error())

}
func TestResetPassword_Success(t *testing.T) {
	userStore, _, _, service := setup(t)

	hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fail()
	}
	userStore.On("GetUserByID", mock.Anything, mock.Anything).Return(&model.User{Password: string(hash)}, nil)
	userStore.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)

	err = service.ResetPassword(t.Context(), uuid.New(), dto.ResetUserPasswordDTO{OldPassword: "123"})

	assert.NoError(t, err)

}
