package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeController struct {
	Storage *database.Storage
}

func NewLikeController(storage *database.Storage) *LikeController {
	return &LikeController{
		Storage: storage,
	}
}

// LikePost godoc
//
//	@Summary		Like a post
//	@Description	Like a post with post ID
//	@Tags			PostLikes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router	/posts/{id}/like [post]
func (lc LikeController) LikePost(c *gin.Context) {

	id := c.Param("id")

	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	liked, err := lc.Storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if liked {
		c.JSON(400, util.ErrorResponse{Error: "Post already liked"})
		return
	}

	err = lc.Storage.LikeStore.LikePost(&model.PostLike{
		PostID: postID,
		UserID: userID,
	})
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	c.JSON(201, util.SuccessMessageResponse{Message: "Post liked successfully"})

}

// UnlikePost godoc
//
//	@Summary		Unlike a post
//	@Description	Unlike a post with post ID
//	@Tags			PostLikes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router	/posts/{id}/unlike [delete]
func (lc LikeController) UnlikePost(c *gin.Context) {
	id := c.Param("id")

	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	liked, err := lc.Storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if !liked {
		c.JSON(400, util.ErrorResponse{Error: "Post not liked yet"})
		return
	}

	err = lc.Storage.LikeStore.UnlikePost(postID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Post unliked successfully"})
}

// LikePost godoc
//
//	@Summary		Like a Comment
//	@Description	Like a Comment with Comment ID
//	@Tags			CommentLikes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//
// @Success		200	{object}	util.SuccessMessageResponse
// @Failure		400	{object}	util.ErrorResponse
// @Failure		401	{object}	util.ErrorResponse
// @Failure		404	{object}	util.ErrorResponse
// @Failure		500	{object}	util.ErrorResponse
// @Security		Bearer
// @Router	/comments/{id}/like [post]
func (lc *LikeController) LikeComment(c *gin.Context) {
	id := c.Param("id")

	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	existLike, err := lc.Storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if existLike {
		c.JSON(400, util.ErrorResponse{Error: "Comment already liked"})
		return

	}
	err = lc.Storage.LikeStore.LikeComment(&model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Comment liked succesfully"})

}

// UnlikeComment godoc
//
//	@Summary		Unlike a comment
//	@Description	Unlike a comment with comment ID
//	@Tags			CommentLikes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router	/comments/{id}/unlike [delete]
func (lc *LikeController) UnlikeComment(c *gin.Context) {
	id := c.Param("id")

	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	existLike, err := lc.Storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if !existLike {
		c.JSON(400, util.ErrorResponse{Error: "Comment not liked yet"})
		return
	}

	err = lc.Storage.LikeStore.UnlikeComment(commentID, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Comment unliked succesfully"})

}
