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

func (rc *ReplyController) ReplyComment(c *gin.Context) {
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
		c.JSON(400, gin.H{"error": util.CommentNotFoundError})
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

	c.JSON(201, gin.H{"message": "Reply created successfully"})

}

// TODO: update reply
