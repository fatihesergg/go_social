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
		c.JSON(500, gin.H{"error": "Internal server error"})
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
