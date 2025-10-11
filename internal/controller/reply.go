package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReplyController struct {
	Storage database.Storage
}

func NewReplyController(storage *database.Storage) *ReplyController {
	return &ReplyController{
		Storage: *storage,
	}
}

// GetCommentReplies godoc
//
//	@Summary		Get replies of a comment
//	@Description	Get replies of a comment
//	@Tags			Reply
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/replies/{id} [GET]
//	@Security		Bearer
func (rc *ReplyController) GetCommentReplies(c *gin.Context) {
	id := c.Param("id")

	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	existComment, err := rc.Storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if existComment == nil {
		c.JSON(404, util.SuccessMessageResponse{Message: "Comment not found"})
		return
	}

	replies, err := rc.Storage.ReplyStore.GetRepliesByCommentID(existComment.ID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if replies == nil {
		c.JSON(404, util.SuccessMessageResponse{Message: "Replies not found"})
		return
	}

	result := dto.NewReplyResponse(replies)
	c.JSON(200, util.SuccessResultResponse{Message: "Replies fetched successfully", Result: result})

}

// ReplyCommnt godoc
//
//	@Summary		Reply a comment
//	@Description	Reply a comment
//	@Tags			Reply
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		201	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/comments/{id}/reply [POST]
//	@Security		Bearer
func (rc *ReplyController) ReplyComment(c *gin.Context) {
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

	var params dto.CreateReply
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	comment, err := rc.Storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		c.JSON(500, util.InternalServerError)
		return
	}

	if comment == nil {
		c.JSON(400, util.ErrorResponse{Error: util.CommentNotFoundError})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	reply := &model.Reply{
		CommentID: comment.ID,
		UserID:    userID,
		Message:   params.Message,
	}

	err = rc.Storage.ReplyStore.CreateReply(reply)
	if err != nil {
		c.JSON(500, util.InternalServerError)
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Reply created successfully"})

}

// UpdateReply godoc
//
//		@Summary		Update a reply
//		@Description	Update a reply
//		@Tags			Reply
//		@Accept			json
//		@Produce		json
//		@Param			id	path		string	true	"Comment ID"
//	 @Param 			reply body dto.UpdateReply true "Update reply"
//		@Success		200 {object}	util.SuccessMessageResponse
//		@Failure		400	{object}	util.ErrorResponse
//		@Failure		404	{object}	util.ErrorResponse
//		@Failure		403	{object}	util.ErrorResponse
//		@Failure		500	{object}	util.ErrorResponse
//		@Router			/replies/{id} [PUT]
//		@Security		Bearer
func (rc *ReplyController) UpdateReply(c *gin.Context) {

	id := c.Param("id")

	var params dto.UpdateReply

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	replyID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}
	existReply, err := rc.Storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if existReply == nil {
		c.JSON(404, util.SuccessMessageResponse{Message: "Reply not found"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	if existReply.UserID != userID {
		c.JSON(403, util.ErrorResponse{Error: util.InvalidPermissionError})
		return
	}

	existReply.Message = params.Message

	err = rc.Storage.ReplyStore.UpdateReply(existReply)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Reply updated successfully"})

}

// DeleteReply godoc
//
//	@Summary		Delete a reply
//	@Description	Delete a reply
//	@Tags			Reply
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200 {object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		403	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/replies/{id} [DELETE]
//	@Security		Bearer
func (rc *ReplyController) DeleteReply(c *gin.Context) {

	id := c.Param("id")

	replyID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	existReply, err := rc.Storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	if existReply == nil {
		c.JSON(404, util.ErrorResponse{Error: "Reply not found"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	if existReply.UserID != userID {
		c.JSON(403, util.ErrorResponse{Error: util.InvalidPermissionError})
		return
	}

	err = rc.Storage.ReplyStore.DeleteReply(existReply.ID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	c.JSON(200, util.SuccessMessageResponse{Message: "Reply deleted successfully"})
}
