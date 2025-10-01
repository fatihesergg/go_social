package controller

import (
	"strconv"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/gin-gonic/gin"
)

type LikeController struct {
	Storage database.Storage
}

func NewLikeController(storage database.Storage) *LikeController {
	return &LikeController{
		Storage: storage,
	}
}

func (lc LikeController) LikePost(c *gin.Context) {

	var params struct {
		PostID int64 `json:"post_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	userID := c.MustGet("userID").(int)

	liked, err := lc.Storage.LikeStore.IsLiked(params.PostID, int64(userID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if liked {
		c.JSON(400, gin.H{"error": "Post already liked"})
		return
	}

	err = lc.Storage.LikeStore.CreateLike(model.Like{
		PostID: params.PostID,
		UserID: int64(userID),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(201, gin.H{"message": "Post liked successfully"})

}

func (lc LikeController) UnlikePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}
	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Post ID"})
		return
	}

	userID := c.MustGet("userID").(int)

	liked, err := lc.Storage.LikeStore.IsLiked(postID, int64(userID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if !liked {
		c.JSON(400, gin.H{"error": "Post not liked yet"})
		return
	}

	err = lc.Storage.LikeStore.DeleteLike(postID, int64(userID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Post unliked successfully"})
}
