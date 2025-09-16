package controller

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/gin-gonic/gin"
)

type FeedController struct {
	Storage database.Storage
}

func NewFeedController(storage database.Storage) FeedController {
	return FeedController{
		Storage: storage,
	}
}

func (fc FeedController) GetFeed(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	search := database.NewSearch(c)
	pagination := database.NewPagination(c)
	posts, err := fc.Storage.FeedStore.GetFeed(int64(userID), pagination, search)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "No posts found"})
			return
		}
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, posts)
}
