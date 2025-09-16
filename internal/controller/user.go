package controller

import (
	"fmt"
	"strconv"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Storage database.Storage
}

func (uc UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": util.UserNotFoundError})
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
//	@Param			user	body		model.User												true	"User signup data"
//	@Success		201		{object}	util.SuccessResponse{result=model.User,message=string}	"User registered successfully"
//	@Failure		400		{object}	util.ErrorResponse{error=string}						"Bad Request: Invalid input or email/username already exists"
//	@Failure		500		{object}	util.ErrorResponse{error=string}						"Internal Server Error"
//	@Router			/signup [post]
func (uc UserController) Signup(c *gin.Context) {
	var params struct {
		Name     string  `json:"name" binding:"required"`
		LastName string  `json:"last_name" binding:"required"`
		Email    string  `json:"email" binding:"required,email"`
		Avatar   *string `json:"avatar"`
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Name:     params.Name,
		LastName: params.LastName,
		Email:    params.Email,
		Avatar:   params.Avatar,
		Username: params.Username,
		Password: params.Password,
	}

	existEmail, err := uc.Storage.UserStore.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting user by email")
		return
	}
	if existEmail != nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		fmt.Println(err)
		return
	}
	existUsername, err := uc.Storage.UserStore.GetUserByUsername(user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		fmt.Println(err)
		return
	}
	if existUsername != nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Something went wrong"})
		return
	}
	user.Password = string(hashedPass)

	err = uc.Storage.UserStore.CreateUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(201, gin.H{"result": user, "message": "User registered successfully"})
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a JWT token
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		object{email=string,password=string}				true	"User login credentials"
//	@Success		200			{object}	util.SuccessResponse{result=string,message=string}	"Login successful"
//	@Failure		400			{object}	util.ErrorResponse{error=string}					"Bad Request: Invalid input"
//	@Failure		401			{object}	util.ErrorResponse{error=string}					"Unauthorized: Invalid credentials"
//	@Failure		404			{object}	util.ErrorResponse{error=string}					"Not Found: User not found"
//	@Failure		500			{object}	util.ErrorResponse{error=string}					"Internal Server Error"
//	@Router			/login [post]
func (uc UserController) Login(c *gin.Context) {
	var params struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByEmail(params.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": util.UserNotFoundError})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	c.JSON(200, gin.H{"result": token, "message": "Login successful"})
}

// GetMe godoc
//
//	@Summary		Get current user
//	@Description	Retrieve the authenticated user's details
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.User
//	@Failure		401	{object}	util.ErrorResponse	"Unauthorized: Invalid or missing token"
//	@Failure		404	{object}	util.ErrorResponse	"Not Found: User not found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Security		Bearer
//	@Router			/me [get]
func (uc UserController) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(int)
	intID := int64(id)
	user, err := uc.Storage.UserStore.GetUserByID(intID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": util.UserNotFoundError})
		return
	}

	c.JSON(200, gin.H{"result": user})
}

// GetFollowerByUserID godoc
//
//	@Summary		Get followers of a user by user ID
//	@Description	Retrieve a list of followers for a specific user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		model.User
//	@Failure		400	{object}	util.ErrorResponse	"Bad Request: ID is required"
//	@Failure		404	{object}	util.ErrorResponse	"Not Found: No followers found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Security		Bearer
//	@Router			/users/{id}/followers [get]
func (uc UserController) GetFollowerByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if len(followers) == 0 {
		c.JSON(404, gin.H{"error": "No followers found"})
		return
	}

	c.JSON(200, gin.H{"result": followers})
}

// GetFollowingByUserID godoc
//
//	@Summary		Get followings of a user by user ID
//	@Description	Retrieve a list of users that a specific user is following
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		model.User
//	@Failure		400	{object}	util.ErrorResponse	"Bad Request: ID is required"
//	@Failure		404	{object}	util.ErrorResponse	"Not Found: No followings found"
//	@Failure		500	{object}	util.ErrorResponse	"Internal Server Error"
//	@Security		Bearer
//	@Router			/users/{id}/following [get]
func (uc UserController) GetFollowingByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	followings, err := uc.Storage.FollowStore.GetFollowingByUserID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if len(followings) == 0 {
		c.JSON(404, gin.H{"error": "No followings found"})
		return
	}

	c.JSON(200, gin.H{"result": followings})
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int															true	"User ID to follow"
//	@Success		200	{object}	util.SuccessResponse{result=model.Follow,message=string}	"Followed successfully"
//	@Failure		400	{object}	util.ErrorResponse											"Bad Request: ID is required"
//	@Failure		500	{object}	util.ErrorResponse											"Internal Server Error"
//	@Security		Bearer
//	@Router			/users/{id}/follow [post]
func (uc UserController) FollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	model := model.Follow{
		UserID:   idInt,
		FollowID: int64(c.MustGet("userID").(int)),
	}

	err = uc.Storage.FollowStore.FollowUser(model.UserID, model.FollowID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Error following user"})
		return
	}
	c.JSON(200, gin.H{"result": model, "message": "Followed successfully"})
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int															true	"User ID to unfollow"
//	@Success		200	{object}	util.SuccessResponse{result=model.Follow,message=string}	"Unfollowed successfully"
//	@Failure		400	{object}	util.ErrorResponse											"Bad Request: ID is required"
//	@Failure		500	{object}	util.ErrorResponse											"Internal Server Error"
//	@Security		Bearer
//	@Router			/users/{id}/unfollow [post]
func (uc UserController) UnfollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	model := model.Follow{
		UserID:   idInt,
		FollowID: int64(c.MustGet("userID").(int)),
	}
	err = uc.Storage.FollowStore.UnFollowUser(model.UserID, model.FollowID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error unfollowing user"})
		return
	}
	c.JSON(200, gin.H{"result": model, "message": "Unfollowed successfully"})
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
//	@Success		200		{array}		model.Post
//	@Failure		400		{object}	util.ErrorResponse	"Bad Request: Invalid user ID"
//	@Failure		403		{object}	util.ErrorResponse	"Forbidden: Not following the user"
//	@Failure		404		{object}	util.ErrorResponse	"Not Found: User not found or no posts found"
//	@Failure		500		{object}	util.ErrorResponse	"Internal Server Error"
//	@Security		Bearer
//	@Router			/users/{id}/posts [get]
func (uc UserController) GetUsersPosts(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(int64(idInt))
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": util.UserNotFoundError})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	userID := c.MustGet("userID").(int)

	pagination := database.NewPagination(c)
	search := database.NewSearch(c)
	for i := range followers {
		followerID := followers[i].ID
		if followerID == int64(userID) {
			posts, err := uc.Storage.PostStore.GetPostsByUserID(user.ID, pagination, search)
			if err != nil {
				c.JSON(500, gin.H{"error": util.InternalServerError})
				return
			}
			if posts == nil {
				c.JSON(404, gin.H{"error": util.NoPostsFoundError})
				return
			}

			c.JSON(200, gin.H{"result": posts})
			return
		}
	}

	c.JSON(403, gin.H{"error": util.NotFollowingError})
}
