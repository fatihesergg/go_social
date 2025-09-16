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
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{"result": user})
}

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
		c.JSON(500, gin.H{"error": "Internal server error"})
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
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := util.CreateJsonWebToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"result": token, "message": "Login successful"})
}

func (uc UserController) GetMe(c *gin.Context) {
	id := c.MustGet("userID").(int)
	intID := int64(id)
	user, err := uc.Storage.UserStore.GetUserByID(intID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{"result": user})
}

func (uc UserController) GetFollowerByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if len(followers) == 0 {
		c.JSON(404, gin.H{"error": "No followers found"})
		return
	}

	c.JSON(200, gin.H{"result": followers})
}

func (uc UserController) GetFollowingByUserID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	followings, err := uc.Storage.FollowStore.GetFollowingByUserID(idInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if len(followings) == 0 {
		c.JSON(404, gin.H{"error": "No followings found"})
		return
	}

	c.JSON(200, gin.H{"result": followings})
}

func (uc UserController) FollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
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

func (uc UserController) UnfollowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
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

func (uc UserController) GetUsersPosts(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.Storage.UserStore.GetUserByID(int64(idInt))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	followers, err := uc.Storage.FollowStore.GetFollowerByUserID(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
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
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			if posts == nil {
				c.JSON(404, gin.H{"error": "No posts found"})
				return
			}

			c.JSON(200, gin.H{"result": posts})
			return
		}
	}

	c.JSON(403, gin.H{"error": "You are not allowed to view this user's posts"})
}
