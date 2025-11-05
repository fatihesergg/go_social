package controller

import (
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController struct {
	PostService services.BasePostService
}

func NewPostController(postService services.BasePostService) *PostController {
	return &PostController{
		PostService: postService,
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
	limit := c.Query("limit")
	offset := c.Query("offset")
	searchQuery := c.Query("search")
	userID := c.MustGet("userID").(uuid.UUID)

	posts, err := pc.PostService.GetAllPosts(userID, limit, offset, searchQuery)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResultResponse{Message: "Posts fetched successfully", Result: posts})
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
	meID := c.MustGet("userID").(uuid.UUID)

	post, err := pc.PostService.GetPostByID(meID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(200, util.SuccessResultResponse{Message: "Post fetched successfully", Result: post})
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
	userID := c.MustGet("userID").(uuid.UUID)

	err := pc.PostService.CreatePost(userID, params)
	if err != nil {
		util.WriteAppError(c, err)
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
//	@Param			id		path		int					true	"Post ID"
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
	userID := c.MustGet("userID").(uuid.UUID)

	err := pc.PostService.UpdatePost(userID, id, params)
	if err != nil {
		util.WriteAppError(c, err)
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

	userID := c.MustGet("userID").(uuid.UUID)
	err := pc.PostService.DeletePost(userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessMessageResponse{Message: "Post deleted successfully"})
}
