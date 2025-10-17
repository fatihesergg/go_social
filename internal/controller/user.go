package controller

import (
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService services.BaseUserService
}

func NewUserController(userService services.BaseUserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// TODO: Add  docs
func (uc UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}

	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, gin.H{"result": user})
}

// Signup godoc
//
//	@Summary		User signup
//	@Description	Register a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateUserDTO												true	"User signup data"
//	@Success		201		{object}	util.SuccessMessageResponse{result=model.User}	"User registered successfully"
//	@Failure		400		{object}	util.ErrorResponse{}
//	@Failure		500		{object}	util.ErrorResponse{}
//	@Router			/signup [post]
func (uc UserController) Signup(c *gin.Context) {
	var params dto.CreateUserDTO

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	err := uc.UserService.Register(params)

	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "User registered successfully"})
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a JWT token
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		dto.LoginUserDTO				true	"User login credentials"
//	@Success		200			{object}	util.SuccessMessageResponse{result=string}
//	@Failure		400			{object}	util.ErrorResponse{error=string}
//	@Failure		401			{object}	util.ErrorResponse{error=string}
//	@Failure		404			{object}	util.ErrorResponse{error=string}
//	@Failure		500			{object}	util.ErrorResponse{error=string}
//	@Router			/login [post]
func (uc UserController) Login(c *gin.Context) {
	var params dto.LoginUserDTO

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	token, err := uc.UserService.Login(params)

	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "Login successful", Result: token})
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Retrieve the authenticated user's details
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	util.SuccessResultResponse{result=model.User}
//	@Failure		401	{object}	util.ErrorResponse	"Unauthorized: Invalid or missing token"
//	@Failure		404	{object}	util.ErrorResponse	"Not Found: User not found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Security		Bearer
//	@Router			/me [get]
func (uc UserController) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(uuid.UUID)
	user, err := uc.UserService.GetUserByID(id.String())
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "User fetched successfully", Result: user})
}

// GetFollowerByUserID godoc
//
//	@Summary		Get followers of a user by user ID
//	@Description	Retrieve a list of followers for a specific user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}		util.SuccessResultResponse{result=[]model.Follow}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/{id}/followers [get]
func (uc UserController) GetFollowerByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}
	followers, err := uc.UserService.GetFollowerByUserID(id)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "Followers fetched successfully", Result: followers})
}

// GetFollowingByUserID godoc
//
//	@Summary		Get followings of a user by user ID
//	@Description	Retrieve a list of users that a specific user is following
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}		util.SuccessResultResponse{result=[]model.Follow}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/{id}/following [get]
func (uc UserController) GetFollowingByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}

	followings, err := uc.UserService.GetFollowingByUserID(id)

	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "User following fetched successfully", Result: followings})
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int															true	"User ID to follow"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/{id}/follow [post]
func (uc UserController) FollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}
	me := c.MustGet("userID").(uuid.UUID)

	err := uc.UserService.FollowUser(me, id)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Followed successfully"})
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int															true	"User ID to unfollow"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/{id}/unfollow [post]
func (uc UserController) UnfollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}
	me := c.MustGet("userID").(uuid.UUID)
	err := uc.UserService.UnFollowUser(me, id)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Unfollowed successfully"})
}

// GetUsersPosts godoc
//
//	@Summary		Get posts of a user by user ID
//	@Description	Retrieve posts made by a specific user, only if the requester is following that user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"User ID"
//	@Param			limit	query		int		false	"Limit"		default(10)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Param			search	query		string	false	"Search query"
//	@Success		200		{object}		util.SuccessResultResponse{result=[]dto.AllPostResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		403		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/{id}/posts [get]
func (uc UserController) GetUsersPosts(c *gin.Context) {
	id := c.Param("id")
	meID := c.MustGet("userID").(uuid.UUID)

	limit := c.Query("limit")
	offset := c.Query("offset")
	searchQuery := c.Query("search")

	posts, err := uc.UserService.GetUsersPosts(id, meID, limit, offset, searchQuery)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Posts fetched succesfully", Result: posts})
}

// ResetPassword godoc
//
//	@Summary		Reset user password
//	@Description	Allow authenticated users to reset their password by providing the old and new passwords
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			passwords	body		dto.ResetUserPasswordDTO	true	"Old and new passwords"
//	@Success		200			{object}	util.SuccessMessageResponse
//	@Failure		400			{object}	util.ErrorResponse
//	@Failure		404			{object}	util.ErrorResponse
//	@Failure		500			{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/users/reset_password [post]
func (uc UserController) ResetPassword(c *gin.Context) {

	meID := c.MustGet("userID").(uuid.UUID)
	var params dto.ResetUserPasswordDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	err := uc.UserService.ResetPassword(meID, params)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Password updated successfully"})
}

//TODO:  Add docs

func (uc UserController) SearchUserByUsername(c *gin.Context) {
	username := c.Param("username")

	users, err := uc.UserService.GetUsersByUsername(username)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "Users fetched successfully", Result: users})

}
