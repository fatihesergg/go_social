package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Storage *database.Storage
}

func NewUserController(storage *database.Storage) *UserController {
	return &UserController{
		Storage: storage,
	}
}

func (uc UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: "ID is required"})
		return
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, util.ErrorResponse{Error: util.UserNotFoundError})
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

	user := &model.User{
		Name:     params.Name,
		LastName: params.LastName,
		Email:    params.Email,
		Avatar:   params.Avatar,
		Username: params.Username,
		Password: params.Password,
	}

	existEmail, err := uc.Storage.UserStore.GetUserByEmail(user.Email)
	if err != nil {

		return
	}
	if existEmail != nil {
		c.JSON(400, util.ErrorResponse{Error: "Email already exists"})

		return
	}
	existUsername, err := uc.Storage.UserStore.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})

		return
	}
	if existUsername != nil {
		c.JSON(400, util.ErrorResponse{Error: "Username already exists"})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Something went wrong"})
		return
	}
	user.Password = string(hashedPass)

	err = uc.Storage.UserStore.CreateUser(user)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error creating user"})
		return
	}

	c.JSON(201, util.SuccessResultResponse{Message: "User registered successfully", Result: user})
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

	user, err := uc.Storage.UserStore.GetUserByEmail(params.Email)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if user == nil {
		c.JSON(404, util.ErrorResponse{Error: util.UserNotFoundError})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		c.JSON(401, util.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
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
	user, err := uc.Storage.UserStore.GetUserByID(id)
	if err != nil {

		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, util.ErrorResponse{Error: util.UserNotFoundError})
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
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if len(followers) == 0 {
		c.JSON(404, util.ErrorResponse{Error: "No followers found"})
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

	userID, err := uuid.Parse(id)

	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	followings, err := uc.Storage.FollowStore.GetFollowingByUserID(userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if len(followings) == 0 {
		c.JSON(404, util.ErrorResponse{Error: "No followings found"})
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
	followUser, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	me := c.MustGet("userID").(uuid.UUID)

	followings, err := uc.Storage.FollowStore.GetFollowingByUserID(me)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	isFollowing := false
	for _, follow := range followings {
		if follow.FollowID == followUser {
			isFollowing = true
		}
	}

	if isFollowing {
		c.JSON(400, util.ErrorResponse{Error: "You are already following this user"})
		return
	}

	err = uc.Storage.FollowStore.FollowUser(me, followUser)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error following user"})
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
	unfUser, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}
	me := c.MustGet("userID").(uuid.UUID)

	followings, err := uc.Storage.FollowStore.GetFollowingByUserID(me)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	isFollowing := false
	for _, follow := range followings {
		if follow.FollowID == unfUser {
			isFollowing = true
		}
	}
	if !isFollowing {
		c.JSON(400, util.ErrorResponse{Error: "You are not following this user"})
		return
	}

	err = uc.Storage.FollowStore.UnFollowUser(me, unfUser)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error unfollowing user"})
		return
	} else {
		c.JSON(200, util.SuccessMessageResponse{Message: "Unfollowed successfully"})
	}
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

	userID, err := uuid.Parse(id)

	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, util.ErrorResponse{Error: util.UserNotFoundError})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(user.ID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	pagination := database.NewPagination(c)
	search := database.NewSearch(c)
	for i := range followers {
		//TODO:  Check if user owner the
		followerID := followers[i].ID
		if followerID == userID {
			posts, err := uc.Storage.PostStore.GetPostsByUserID(user.ID, pagination, search)
			if err != nil {
				c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
				return
			}
			if posts == nil {
				c.JSON(404, util.ErrorResponse{Error: util.NoPostsFoundError})
				return
			}
			result := dto.NewAllPostResponse(posts)

			c.JSON(200, util.SuccessResultResponse{Message: "User posts fetched successfully", Result: result})
			return
		}
	}

	c.JSON(403, util.ErrorResponse{Error: util.NotFollowingError})
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

	var params dto.ResetUserPasswordDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	user, err := uc.Storage.UserStore.GetUserByID(userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if user == nil {
		c.JSON(404, util.ErrorResponse{Error: util.UserNotFoundError})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.OldPassword)); err != nil {
		c.JSON(400, util.ErrorResponse{Error: "Invalid Credentials"})
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Something went wrong"})
		return
	}
	user.Password = string(hashedPass)

	err = uc.Storage.UserStore.UpdateUser(user)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error updating password"})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Password updated successfully"})
}

func (uc UserController) SearchUserByUsername(c *gin.Context) {
	username := c.Param("username")

	users, err := uc.Storage.UserStore.GetUsersByUsername(username)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if users == nil {
		c.JSON(200, util.SuccessMessageResponse{Message: "No users found"})
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "Users fetched successfully", Result: users})

}
