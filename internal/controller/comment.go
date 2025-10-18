package controller

import (
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	CommentService services.BaseCommentService
}

func NewCommentController(commentService services.BaseCommentService) *CommentController {
	return &CommentController{
		CommentService: commentService,
	}
}

// CreateComment godoc
//
//	@Summary		Create a new comment
//	@Description	Create a new comment on a post
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			comment	body		dto.CreateCommentDTO	true	"Comment to create"
//	@Success		201		{object}	util.SuccessMessageResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	var params dto.CreateCommentDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)
	err := cc.CommentService.AddCommentPost(userID, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Comment created successfully"})

}

// GetCommentsByPostID godoc
//
//	@Summary		Get comments for a specific post
//	@Description	Retrieve all comments associated with a specific post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200		{object}	util.SuccessResultResponse{result=[]dto.CommentDetailResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/posts/{id}/comments [get]
func (cc *CommentController) GetCommentsByPostID(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)

	comments, err := cc.CommentService.GetCommentsByPostID(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Comments fetched successfully", Result: comments})
}

// UpdateComment godoc
//
//	@Summary		Update a comment
//	@Description	Update an existing comment by its ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Comment ID"
//	@Param			comment	body		dto.UpdateCommentDTO	true	"Updated comment data"
//	@Success		200		{object}	util.SuccessMessageResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/{id} [put]
func (cc *CommentController) UpdateComment(c *gin.Context) {
	var params dto.UpdateCommentDTO
	id := c.Param("id")

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	err := cc.CommentService.UpdateComment(userID, id, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Comment updated successfully"})
}

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete an existing comment by its ID
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
//	@Router	/comments/{id} [delete]
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)

	err := cc.CommentService.DeleteComment(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Comment deleted successfully"})
}
