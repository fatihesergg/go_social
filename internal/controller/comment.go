package controller

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
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

func (cc CommentController) CreateComment(c *gin.Context) {
	var params struct {
		PostID  int64  `json:"post_id" binding:"required"`
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	userID := c.MustGet("userID").(int)
	comment := model.Comment{
		ID:      uuid.New(),
		PostID:  params.PostID,
		UserID:  int64(userID),
		Content: params.Content,
	}
	fmt.Println(comment)
	if params.Image != "" {
		comment.Image = sql.NullString{String: params.Image, Valid: true}
	} else {
		comment.Image = sql.NullString{}
	}

	err := cc.Storage.CommentStore.CreateComment(comment)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Error creating comment"})
		return
	}
	c.JSON(201, gin.H{"result": "", "message": "Comment created successfully"})

}

func (cc CommentController) GetCommentsByPostID(c *gin.Context) {
	id := c.Param("post_id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	comments, err := cc.Storage.CommentStore.GetCommentsByPostID(int64(intID))
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if comments == nil {
		c.JSON(404, gin.H{"error": "No comments found"})
		return
	}
	c.JSON(200, gin.H{"result": comments})
}

func (cc CommentController) UpdateComment(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
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
		c.JSON(404, gin.H{"error": "Comment not found"})
		return
	}
	comment.Content = params.Content
	if params.Image != "" {
		comment.Image = sql.NullString{String: params.Image, Valid: true}
	} else {
		comment.Image = sql.NullString{}
	}
	err = cc.Storage.CommentStore.UpdateComment(*comment)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating comment"})
		return
	}
	c.JSON(200, gin.H{"result": "", "message": "Comment updated successfully"})
}

func (cc CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is required"})
		return
	}
	commentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	err = cc.Storage.CommentStore.DeleteComment(commentID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting comment"})
		return
	}
	c.JSON(200, gin.H{"result": "", "message": "Comment deleted successfully"})
}
