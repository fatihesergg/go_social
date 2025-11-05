package controller

import (
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeController struct {
	LikeService services.BaseLikeService
}

func NewLikeController(likeService services.BaseLikeService) *LikeController {
	return &LikeController{
		LikeService: likeService,
	}
}

// LikePost godoc
//
//	@Summary		Like a post
//	@Description	Like a post with post ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/posts/{id}/like [post]
func (lc LikeController) LikePost(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	err := lc.LikeService.LikePost(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(201, util.SuccessMessageResponse{Message: "Post liked successfully"})

}

// LikeReply godoc
//
//	@Summary		Like a reply
//	@Description	Like a reply with reply ID
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Reply ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/replies/{id}/like [post]
func (lc LikeController) LikeReply(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	err := lc.LikeService.LikeReply(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(201, util.SuccessMessageResponse{Message: "Reply liked successfully"})

}

// UnlikePost godoc
//
//	@Summary		Unlike a post
//	@Description	Unlike a post with post ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/posts/{id}/unlike [delete]
func (lc LikeController) UnlikePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	err := lc.LikeService.UnlikePost(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Post unliked successfully"})
}

// UnlikeReply godoc
//
//	@Summary		Unlike a reply
//	@Description	Unlike a reply with reply ID
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Reply ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/replies/{id}/unlike [delete]
func (lc LikeController) UnlikeReply(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	err := lc.LikeService.UnlikeReply(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Reply unliked successfully"})
}

// LikePost godoc
//
//	@Summary		Like a Comment
//	@Description	Like a Comment with Comment ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/{id}/like [post]
func (lc *LikeController) LikeComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	err := lc.LikeService.LikeComment(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(201, util.SuccessMessageResponse{Message: "Comment liked succesfully"})

}

// UnlikeComment godoc
//
//	@Summary		Unlike a comment
//	@Description	Unlike a comment with comment ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/{id}/unlike [delete]
func (lc *LikeController) UnlikeComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)

	err := lc.LikeService.UnlikeComment(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Comment unliked succesfully"})

}
