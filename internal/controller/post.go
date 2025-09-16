package controller

import (
	"fmt"
	"strconv"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	Storage database.Storage
}

func NewPostController(storage database.Storage) *PostController {
	return &PostController{
		Storage: storage,
	}
}

// GetPosts godoc
//
//	@Summary		Get all posts
//	@Description	Retrieve a list of all posts with optional pagination and search
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Param			search	query		string	false	"Search query"
//	@Success		200		{array}		model.Post
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts [get]
//	@Security		Bearer
func (pc PostController) GetPosts(c *gin.Context) {
	pagination := database.NewPagination(c)
	search := database.NewSearch(c)
	posts, err := pc.Storage.PostStore.GetPosts(pagination, search)
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if posts == nil {
		c.JSON(404, gin.H{"error": util.NoPostsFoundError})
		return
	}
	c.JSON(200, gin.H{"result": posts})
}

// GetPostByID godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a single post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	model.Post
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/posts/{id} [get]
//	@Security		Bearer
func (pc PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}
	post, err := pc.Storage.PostStore.GetPostByID(int64(intID))
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if post == nil {
		c.JSON(404, gin.H{"error": util.PostNotFoundError})
		return
	}

	c.JSON(200, gin.H{"result": post})
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with content and optional image
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		model.Post	true	"Post data"
//	@Success		201		{object}	model.Post
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts [post]
//	@Security		Bearer
func (pc PostController) CreatePost(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	post := model.Post{
		Content: params.Content,
	}
	if params.Image != "" {
		post.Image.String = params.Image
	}
	userID := c.MustGet("userID").(int)

	post.UserID = int64(userID)

	err := pc.Storage.PostStore.CreatePost(post)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}

	c.JSON(201, gin.H{"result": post})

}

// UpdatePost godoc
//
//	@Summary		Update an existing post
//	@Description	Update the content and/or image of an existing post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Post ID"
//	@Param			post	body		model.Post	true	"Updated post data"
//	@Success		200		{object}	model.Post
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		403		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts/{id} [put]
//	@Security		Bearer
func (pc PostController) UpdatePost(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		Image   string `json:"image"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}

	existPost, err := pc.Storage.PostStore.GetPostByID(int64(intID))
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if existPost == nil {
		c.JSON(404, gin.H{"error": util.PostNotFoundError})
		return
	}
	if existPost.UserID != int64(c.MustGet("userID").(int)) {
		c.JSON(403, gin.H{"error": util.UnauthorizedError})
		return
	}
	post := model.Post{
		ID:      int64(intID),
		Content: params.Content,
	}

	if params.Image != "" {
		post.Image.String = params.Image
	} else {
		post.Image.Valid = false
	}

	err = pc.Storage.PostStore.UpdatePost(post)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating post"})
		return
	}
	c.JSON(200, gin.H{"result": post})

}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete an existing post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		403	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/posts/{id} [delete]
//	@Security		Bearer
func (pc PostController) DeletePost(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": util.IDRequiredError})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": util.InvalidIDFormatError})
		return
	}
	post, err := pc.Storage.PostStore.GetPostByID(int64(idInt))
	if err != nil {
		c.JSON(500, gin.H{"error": util.InternalServerError})
		return
	}
	if post == nil {
		c.JSON(404, gin.H{"error": util.PostNotFoundError})
		return
	}
	if post.UserID != int64(c.MustGet("userID").(int)) {
		c.JSON(403, gin.H{"error": util.UnauthorizedError})
		return
	}

	err = pc.Storage.PostStore.DeletePost(int64(idInt))
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting post"})
		return
	}
	c.JSON(200, gin.H{"message": "Post deleted successfully"})
}
