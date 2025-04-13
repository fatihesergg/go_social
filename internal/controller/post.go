package controller

import (
	"fmt"
	"strconv"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	Storage database.Storage
}

func NewPostController(storage database.Storage) *PostController {
	return &PostController{
		Storage: storage,
	}
}

func (pc PostController) GetPosts(c *gin.Context) {
	posts, err := pc.Storage.PostStore.GetPosts()
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if posts == nil {
		c.JSON(404, gin.H{"error": "No posts found"})
		return
	}
	c.JSON(200, gin.H{"result": posts})
}

func (pc PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	post, err := pc.Storage.PostStore.GetPostByID(int64(intID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if post == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, gin.H{"result": post})
}

func (pc PostController) CreatePost(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	post := model.Post{
		Content: params.Content,
	}
	if params.Image != "" {
		post.Image.String = params.Image
	}
	userID := c.MustGet("userID").(int)

	post.UserID = int64(userID)

	err := pc.Storage.PostStore.CreatePost(post)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(201, gin.H{"result": post})

}

func (pc PostController) UpdatePost(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	existPost, err := pc.Storage.PostStore.GetPostByID(int64(intID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if existPost == nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	if existPost.UserID != int64(c.MustGet("userID").(int)) {
		c.JSON(403, gin.H{"error": "You are not authorized to update this post"})
		return
	}
	post := model.Post{
		ID:      int64(intID),
		Content: params.Content,
	}

	if params.Image != "" {
		post.Image.String = params.Image
	} else {
		post.Image.Valid = false
	}

	err = pc.Storage.PostStore.UpdatePost(post)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating post"})
		return
	}
	c.JSON(200, gin.H{"result": post})

}
