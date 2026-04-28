package controller

import (
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReplyController struct {
	ReplyService services.BaseReplyService
}

func NewReplyController(replyService services.BaseReplyService) *ReplyController {
	return &ReplyController{
		ReplyService: replyService,
	}
}

// GetRepliesByParent godoc
//
//	@Summary		Get replies of a reply
//	@Description	Get replies of a reply
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Reply ID"
//	@Success		200	{object}	util.SuccessResultResponse{result=[]dto.ReplyResponse}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/replies/{id}/replies [GET]
//	@Security		Bearer
func (rc *ReplyController) GetRepliesByParent(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	replies, err := rc.ReplyService.GetRepliesByParentID(c.Request.Context(), userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Replies fetched successfully", Result: replies})
}

// GetCommentReplies godoc
//
//	@Summary		Get replies of a comment
//	@Description	Get replies of a comment
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessResultResponse{result=[]dto.ReplyResponse}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/comments/{id}/replies [GET]
//	@Security		Bearer
func (rc *ReplyController) GetCommentReplies(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	replies, err := rc.ReplyService.GetCommentReplies(c.Request.Context(), userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Replies fetched successfully", Result: replies})

}

// ReplyComment godoc
//
//	@Summary		Reply a comment
//	@Description	Reply a comment
//	@Tags			Comments
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string			true	"Comment ID"
//	@Param			CreateReply	body		dto.CreateReply	true	"User login credentials"
//	@Success		201			{object}	util.SuccessMessageResponse
//	@Failure		400			{object}	util.ErrorResponse
//	@Failure		500			{object}	util.ErrorResponse
//	@Router			/comments/{id}/reply [POST]
//	@Security		Bearer
func (rc *ReplyController) ReplyComment(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	var params dto.CreateReply
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	err := rc.ReplyService.ReplyComment(c.Request.Context(), userID, id, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Reply created successfully"})

}

// ReplyComment godoc
//
//	@Summary		Reply a reply
//	@Description	Reply a reply
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string			true	"Reply ID"
//	@Param			CreateReply	body		dto.CreateReply	true	"User login credentials"
//	@Success		201			{object}	util.SuccessMessageResponse
//	@Failure		400			{object}	util.ErrorResponse
//	@Failure		500			{object}	util.ErrorResponse
//	@Router			/replies/{id}/reply [POST]
//	@Security		Bearer
func (rc *ReplyController) ReplyAReply(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	var params dto.CreateReply
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	err := rc.ReplyService.ReplyAReply(c.Request.Context(), userID, id, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Reply created successfully"})

}

// UpdateReply godoc
//
//	@Summary		Update a reply
//	@Description	Update a reply
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"Comment ID"
//	@Param			reply	body		dto.UpdateReply	true	"Update reply"
//	@Success		200		{object}	util.SuccessMessageResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		403		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/replies/{id} [PUT]
//	@Security		Bearer
func (rc *ReplyController) UpdateReply(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)
	var params dto.UpdateReply

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	err := rc.ReplyService.UpdateReply(c.Request.Context(), userID, id, params)

	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Reply updated successfully"})

}

// DeleteReply godoc
//
//	@Summary		Delete a reply
//	@Description	Delete a reply
//	@Tags			Replies
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Comment ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		403	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/replies/{id} [DELETE]
//	@Security		Bearer
func (rc *ReplyController) DeleteReply(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("userID").(uuid.UUID)

	err := rc.ReplyService.DeleteReply(c.Request.Context(), userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Reply deleted successfully"})
}
