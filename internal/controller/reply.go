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

// ReplyCommnt godoc
//
//	@Summary		Reply a comment
//	@Description	Reply a comment
//	@Tags			Reply
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uuid	true	"Comment ID"
//	@Success		201	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/reply/{id} [POST]
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

// TODO: update reply
