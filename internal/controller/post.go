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
	TagService  services.BaseTagService
}

func NewPostController(postService services.BasePostService, tagService services.BaseTagService) *PostController {
	return &PostController{
		PostService: postService,
		TagService:  tagService,
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
//	@Success		200		{array}		util.SuccessResponse{result=[]dto.AllPostResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts [get]
//	@Security		Bearer
func (pc PostController) GetPosts(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	searchQuery := c.Query("search")
	userID := c.MustGet("userID").(uuid.UUID)

	posts, err := pc.PostService.GetAllPosts(c.Request.Context(), userID, limit, offset, searchQuery)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	result := dto.NewAllPostResponse(posts)

	c.JSON(200, util.SuccessResponse{Message: "posts fetched successfully", Result: result})
}

// GetPosts godoc
//
//	@Summary		Get all posts by tag
//	@Description	Retrieve a list of all posts with optional pagination and tag
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Param			tag	path		string	false	"Tag query"
//	@Success		200		{array}		util.SuccessResponse{result=[]dto.AllPostResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/posts/tag/{tag} [get]
//	@Security		Bearer
func (pc PostController) GetPostsByTag(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	tagParam := c.Param("tag")
	userID := c.MustGet("userID").(uuid.UUID)

	posts, err := pc.PostService.GetAllPostsByTag(c.Request.Context(), userID, limit, offset, tagParam)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	result := dto.NewAllPostResponse(posts)
	c.JSON(200, util.SuccessResponse{Message: "posts fetched successfully", Result: result})
}

// GetPostByID godoc
//
//	@Summary		Get a post by ID
//	@Description	Retrieve a single post by its ID
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	util.SuccessResponse{result=dto.PostDetailResponse}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Router			/posts/{id} [get]
//	@Security		Bearer
func (pc PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	meID := c.MustGet("userID").(uuid.UUID)

	post, err := pc.PostService.GetPostByID(c.Request.Context(), meID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	result := dto.NewPostDetailResponse(post)

	c.JSON(200, util.SuccessResponse{Message: "post fetched successfully", Result: result})
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with content and optional image
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		dto.CreatePostDTO	true	"Post data"
//	@Success		201		{object}	util.SuccessResponse
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

	tags, content := pc.TagService.ExtractTagStringFromContent(params.Content)

	params.Content = content

	postID, err := pc.PostService.CreatePost(c.Request.Context(), userID, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	if len(tags) > 0 {
		err = pc.TagService.AddPostTags(c.Request.Context(), postID, tags)
		if err != nil {
			util.WriteAppError(c, err)
			return
		}
	}

	c.JSON(201, util.SuccessResponse{Message: "post created succesfully"})

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
//	@Success		200		{object}	util.SuccessResponse
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

	err := pc.PostService.UpdatePost(c.Request.Context(), userID, id, params)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResponse{Message: "post updated succesfully"})

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

	userID := c.MustGet("userID").(uuid.UUID)
	err := pc.PostService.DeletePost(c.Request.Context(), userID, id)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResponse{Message: "post deleted successfully"})
}
