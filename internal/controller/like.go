package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
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

func (lc LikeController) LikePost(c *gin.Context) {

	var params dto.CreatePostLikeDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	liked, err := lc.Storage.LikeStore.IsPostLiked(params.PostID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if liked {
		c.JSON(400, gin.H{"error": "Post already liked"})
		return
	}

	err = lc.Storage.LikeStore.LikePost(&model.PostLike{
		PostID: params.PostID,
		UserID: userID,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	c.JSON(201, gin.H{"message": "Post liked successfully"})

}

func (lc LikeController) UnlikePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Post ID"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	liked, err := lc.Storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if !liked {
		c.JSON(400, gin.H{"error": "Post not liked yet"})
		return
	}

	err = lc.Storage.LikeStore.UnlikePost(postID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	c.JSON(200, gin.H{"message": "Post unliked successfully"})
}

func (lc *LikeController) LikeComment(c *gin.Context) {
	var params dto.CreateCommentLikeDTO

	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	existLike, err := lc.Storage.LikeStore.IsCommentLiked(params.CommentID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if existLike {
		c.JSON(400, gin.H{"error": "Comment already liked"})
		return

	}
	err = lc.Storage.LikeStore.LikeComment(&model.CommentLike{
		CommentID: params.CommentID,
		UserID:    userID,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	c.JSON(201, gin.H{"message": "Comment liked succesfully"})

}

func (lc *LikeController) UnlikeComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Comment ID required"})
		return
	}

	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Comment ID"})
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	existLike, err := lc.Storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	if !existLike {
		c.JSON(400, gin.H{"error": "Comment not liked yet"})
		return
	}

	err = lc.Storage.LikeStore.UnlikeComment(commentID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	c.JSON(200, gin.H{"message": "Comment unliked succesfully"})

}
