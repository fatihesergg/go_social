package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	Storage *database.Storage
}

func NewCommentController(storage *database.Storage) *CommentController {
	return &CommentController{
		Storage: storage,
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
	comment := &model.Comment{
		ID:      uuid.New(),
		PostID:  params.PostID,
		UserID:  userID,
		Content: params.Content,
	}

	err := cc.Storage.CommentStore.CreateComment(comment)
	if err != nil {

		c.JSON(500, util.ErrorResponse{Error: "Error creating comment"})
		return
	}
	c.JSON(201, util.SuccessMessageResponse{Message: "Comment created successfully"})

}

// GetCommentsByPostID godoc
//
//	@Summary		Get comments for a specific post
//	@Description	Retrieve all comments associated with a specific post by its ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	util.SuccessResultResponse{result=[]dto.CommentDetailResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/post/{post_id} [get]
func (cc *CommentController) GetCommentsByPostID(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: util.IDRequiredError})
		return
	}
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)
	comments, err := cc.Storage.CommentStore.GetCommentsByPostID(postID, userID)

	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if comments == nil {
		c.JSON(404, util.ErrorResponse{Error: util.NoCommentsFoundError})
		return
	}
	result := dto.NewCommentResponse(comments)

	c.JSON(200, util.SuccessResultResponse{Message: "Comments fetched successfully", Result: result})
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
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: util.IDRequiredError})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	comment, err := cc.Storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error fetching comment"})
		return
	}
	if comment == nil {
		c.JSON(404, util.ErrorResponse{Error: util.NoCommentsFoundError})
		return
	}
	comment.Content = params.Content

	err = cc.Storage.CommentStore.UpdateComment(comment)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error updating comment"})
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
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: util.IDRequiredError})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}
	// TODO: CHECK USERID
	err = cc.Storage.CommentStore.DeleteComment(commentID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error deleting comment"})
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Comment deleted successfully"})
}
