package controller

import (
	"github.com/fatihesergg/go_social/internal/errors"
	_ "github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FeedController struct {
	FeedService services.BaseFeedService
}

func NewFeedController(feedService services.BaseFeedService) *FeedController {
	return &FeedController{
		FeedService: feedService,
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
//	@Success		200		{array}		util.SuccessResultResponse{result=[]dto.FeedResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/feed [get]
func (fc FeedController) GetFeed(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	limit := c.Query("limit")
	offset := c.Query("offset")
	searchQuery := c.Query("search")
	posts, err := fc.FeedService.GetFeed(userID, limit, offset, searchQuery)
	if err != nil {
		appErr, _ := err.(errors.AppError)
		c.JSON(appErr.Code, util.ErrorResponse{Error: appErr.Error()})
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Posts fetched successfully", Result: posts})
}
