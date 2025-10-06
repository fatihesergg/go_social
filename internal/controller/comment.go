package controller

import (
	"database/sql"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	Storage database.Storage
}

func NewCommentController(storage database.Storage) *CommentController {
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
//	@Param			comment	body		model.Comment	true	"Comment to create"
//	@Success		201		{object}	util.SuccessResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments [post]
func (cc CommentController) CreateComment(c *gin.Context) {
	var params struct {
		PostID  uuid.UUID `json:"post_id" binding:"required"`
		Content string    `json:"content" binding:"required"`
		Image   string    `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)
	comment := &model.Comment{
		ID:      uuid.New(),
		PostID:  params.PostID,
		UserID:  userID,
		Content: params.Content,
	}

	if params.Image != "" {
		comment.Image = sql.NullString{String: params.Image, Valid: true}
	} else {
		comment.Image = sql.NullString{}
	}

	err := cc.Storage.CommentStore.CreateComment(comment)
	if err != nil {

		c.JSON(500, gin.H{"error": "Error creating comment"})
		return
	}
	c.JSON(201, gin.H{"result": "", "message": "Comment created successfully"})

}

// GetCommentsByPostID godoc
//
//	@Summary		Get comments for a specific post
//	@Description	Retrieve all comments associated with a specific post by its ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	util.SuccessResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/post/{post_id} [get]
func (cc CommentController) GetCommentsByPostID(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}
	comments, err := cc.Storage.CommentStore.GetCommentsByPostID(postID)
	if err != nil {

		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if comments == nil {
		c.JSON(404, gin.H{"error": util.NoCommentsFoundError})
		return
	}
	c.JSON(200, gin.H{"result": comments})
}

// UpdateComment godoc
//
//	@Summary		Update a comment
//	@Description	Update an existing comment by its ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Comment ID"
//	@Param			comment	body		model.Comment	true	"Updated comment data"
//	@Success		200		{object}	util.SuccessResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		401		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/comments/{id} [put]
func (cc CommentController) UpdateComment(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	comment, err := cc.Storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching comment"})
		return
	}
	if comment == nil {
		c.JSON(404, gin.H{"error": util.NoCommentsFoundError})
		return
	}
	comment.Content = params.Content
	if params.Image != "" {
		comment.Image = sql.NullString{String: params.Image, Valid: true}
	} else {
		comment.Image = sql.NullString{}
	}
	err = cc.Storage.CommentStore.UpdateComment(comment)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating comment"})
		return
	}
	c.JSON(200, gin.H{"result": "", "message": "Comment updated successfully"})
}

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete an existing comment by its ID
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		401	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router	/comments/{id} [delete]
func (cc CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}
	err = cc.Storage.CommentStore.DeleteComment(commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting comment"})
		return
	}
	c.JSON(200, gin.H{"result": "", "message": "Comment deleted successfully"})
}
