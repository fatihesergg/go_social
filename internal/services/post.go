package services

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BasePostService interface {
	GetAllPosts(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]dto.AllPostResponse, error)
	GetAllPostsByTag(ctx context.Context, userID uuid.UUID, limit, offset string, tag string) ([]dto.AllPostResponse, error)
	GetPostByID(ctx context.Context, userID uuid.UUID, postIDRaw string) (*dto.PostDetailResponse, error)
	CreatePost(ctx context.Context, userID uuid.UUID, dto dto.CreatePostDTO) (uuid.UUID, error)
	UpdatePost(ctx context.Context, userID uuid.UUID, postIDRaw string, dto dto.UpdatePostDTO) error
	DeletePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error
}
type PostService struct {
	storage *database.Storage
}

func NewPostService(storage *database.Storage) BasePostService {
	return &PostService{storage: storage}
}
func (ps *PostService) GetAllPosts(ctx context.Context, userID uuid.UUID, limit, offset, query string) ([]dto.AllPostResponse, error) {
	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)

	posts, err := ps.storage.PostStore.GetPosts(ctx, pagination, search, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewAllPostResponse(posts)
	return result, nil

}

func (ps *PostService) GetPostByID(ctx context.Context, userID uuid.UUID, postIDRaw string) (*dto.PostDetailResponse, error) {
	postID, err := uuid.Parse(postIDRaw)

	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing postID: %w", err))
	}

	hasAccess, err := ps.storage.PostStore.HasAccessToPost(ctx, userID, postID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this post"))
	}

	post, err := ps.storage.PostStore.GetPostDetailsByID(ctx, postID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("Post not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewPostDetailResponse(post)
	return result, nil
}
func (ps *PostService) CreatePost(ctx context.Context, userID uuid.UUID, dto dto.CreatePostDTO) (uuid.UUID, error) {

	if dto.Visibility == "private" {
		dto.Visibility = "private"
	} else {
		dto.Visibility = "public"
	}
	post := &model.Post{
		ID:         uuid.New(),
		Content:    dto.Content,
		Visibility: model.PostVisibility(dto.Visibility),
	}

	post.UserID = userID

	err := ps.storage.PostStore.CreatePost(ctx, post)
	if err != nil {
		return uuid.Nil, errors.InternalServerError.Wrap(err)
	}

	return post.ID, nil
}
func (ps *PostService) UpdatePost(ctx context.Context, userID uuid.UUID, postIDRaw string, dto dto.UpdatePostDTO) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing postID: %w", err)).Wrap(err)
	}

	existPost, err := ps.storage.PostStore.GetPostByID(ctx, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.PostNotFoundError.Wrap(fmt.Errorf("Post not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existPost.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("Request userid and post userid is different"))
	}
	post := &model.Post{
		ID:      postID,
		Content: dto.Content,
	}

	err = ps.storage.PostStore.UpdatePost(ctx, post)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ps *PostService) DeletePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing postID: %w", err))
	}
	post, err := ps.storage.PostStore.GetPostByID(ctx, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.PostNotFoundError.Wrap(fmt.Errorf("Post not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}
	if post.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("Request userid and post userid is different"))
	}

	err = ps.storage.PostStore.DeletePost(ctx, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ps *PostService) GetAllPostsByTag(ctx context.Context, userID uuid.UUID, limit, offset string, tag string) ([]dto.AllPostResponse, error) {

	pagination := database.NewPagination(limit, offset)
	tagExpr, err := regexp.Compile(`[\w\d]+`)
	if err != nil {
		return nil, errors.InternalServerError
	}
	isValid := tagExpr.MatchString(tag)
	if !isValid {
		return nil, errors.InvalidTag
	}

	posts, err := ps.storage.PostStore.GetPostsByTag(ctx, pagination, tag, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoPostsFoundError.Wrap(fmt.Errorf("No posts found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return dto.NewAllPostResponse(posts), nil

}
