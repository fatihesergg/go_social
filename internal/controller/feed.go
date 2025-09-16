package controller

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/database"
	_ "github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
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

// GetFeed godoc
//
//	@Summary		Get feed posts
//	@Description	Get feed posts for the authenticated user
//	@Tags			Feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Param			search	query		string	false	"Search query"
//	@Success		200		{array}		model.Post
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/feed [get]
func (fc FeedController) GetFeed(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	search := database.NewSearch(c)
	pagination := database.NewPagination(c)
	posts, err := fc.Storage.FeedStore.GetFeed(int64(userID), pagination, search)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": util.PostNotFoundError})
			return
		}
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	c.JSON(200, posts)
}
