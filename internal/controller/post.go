package controller

import (
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController struct {
	Storage *database.Storage
}

func NewPostController(storage *database.Storage) *PostController {
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
//	@Success		200		{array}		util.SuccessResultResponse{result=[]dto.AllPostResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts [get]
//	@Security		Bearer
func (pc PostController) GetPosts(c *gin.Context) {
	pagination := database.NewPagination(c)
	search := database.NewSearch(c)
	userID := c.MustGet("userID").(uuid.UUID)
	posts, err := pc.Storage.PostStore.GetPosts(pagination, search, userID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if posts == nil {
		c.JSON(404, util.ErrorResponse{Error: util.NoPostsFoundError})
		return
	}
	result := dto.NewAllPostResponse(posts)

	c.JSON(200, util.SuccessResultResponse{Message: "Posts fetched successfully", Result: result})
}

// GetPostByID godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a single post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	util.SuccessResultResponse{result=dto.PostDetailResponse}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/posts/{id} [get]
//	@Security		Bearer
func (pc PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
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

	post, err := pc.Storage.PostStore.GetPostDetailsByID(postID, userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if post == nil {
		c.JSON(404, util.ErrorResponse{Error: util.PostNotFoundError})
		return
	}
	result := dto.NewPostDetailResponse(post)

	c.JSON(200, util.SuccessResultResponse{Message: "Post fetched successfully", Result: result})
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with content and optional image
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		dto.CreatePostDTO	true	"Post data"
//	@Success		201		{object}	util.SuccessMessageResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts [post]
//	@Security		Bearer
func (pc PostController) CreatePost(c *gin.Context) {
	var params dto.CreatePostDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}

	post := &model.Post{
		Content: params.Content,
	}

	userID := c.MustGet("userID").(uuid.UUID)
	post.UserID = userID

	err := pc.Storage.PostStore.CreatePost(post)
	if err != nil {

		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}

	c.JSON(201, util.SuccessMessageResponse{Message: "Post created succesfully"})

}

// UpdatePost godoc
//
//	@Summary		Update an existing post
//	@Description	Update the content and/or image of an existing post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Post ID"
//	@Param			post	body		dto.UpdatePostDTO	true	"Updated post data"
//	@Success		200		{object}	util.SuccessMessageResponse
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		403		{object}	util.ErrorResponse
//	@Failure		404		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts/{id} [put]
//	@Security		Bearer
func (pc PostController) UpdatePost(c *gin.Context) {
	var params dto.UpdatePostDTO
	if err := c.ShouldBindJSON(&params); err != nil {
		util.HandleBindError(c, err)
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: util.IDRequiredError})
		return
	}
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}

	existPost, err := pc.Storage.PostStore.GetPostByID(postID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if existPost == nil {
		c.JSON(404, util.ErrorResponse{Error: util.PostNotFoundError})
		return
	}
	if existPost.UserID != c.MustGet("userID").(uuid.UUID) {
		c.JSON(403, util.ErrorResponse{Error: util.UnauthorizedError})
		return
	}
	post := &model.Post{
		ID:      postID,
		Content: params.Content,
	}

	err = pc.Storage.PostStore.UpdatePost(post)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error updating post"})
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Post updated succesfully"})

}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete an existing post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	util.SuccessMessageResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		403	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/posts/{id} [delete]
//	@Security		Bearer
func (pc PostController) DeletePost(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(400, util.ErrorResponse{Error: util.IDRequiredError})
		return
	}
	postID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, util.ErrorResponse{Error: util.InvalidIDFormatError})
		return
	}
	post, err := pc.Storage.PostStore.GetPostByID(postID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: util.InternalServerError})
		return
	}
	if post == nil {
		c.JSON(404, util.ErrorResponse{Error: util.PostNotFoundError})
		return
	}
	if post.UserID != c.MustGet("userID").(uuid.UUID) {
		c.JSON(403, util.ErrorResponse{Error: util.UnauthorizedError})
		return
	}

	err = pc.Storage.PostStore.DeletePost(postID)
	if err != nil {
		c.JSON(500, util.ErrorResponse{Error: "Error deleting post"})
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Post deleted successfully"})
}
