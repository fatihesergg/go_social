package services

import (
	"database/sql"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BasePostService interface {
	GetAllPosts(userID uuid.UUID, limit, offset, query string) ([]dto.AllPostResponse, error)
	GetPostByID(userID uuid.UUID, postIDRaw string) (*dto.PostDetailResponse, error)
	CreatePost(userID uuid.UUID, dto dto.CreatePostDTO) error
	UpdatePost(userID uuid.UUID, postIDRaw string, dto dto.UpdatePostDTO) error
	DeletePost(userID uuid.UUID, postIDRaw string) error
}
type PostService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewPostService(storage *database.Storage, logger *zap.Logger) BasePostService {
	return &PostService{storage: storage, logger: logger.Named("post_service")}
}
func (ps *PostService) GetAllPosts(userID uuid.UUID, limit, offset, query string) ([]dto.AllPostResponse, error) {
	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)

	posts, err := ps.storage.PostStore.GetPosts(pagination, search, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.logger.Error("No posts found")
			return nil, errors.NoPostsFoundError
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewAllPostResponse(posts)
	return result, nil

}

func (ps *PostService) GetPostByID(userID uuid.UUID, postIDRaw string) (*dto.PostDetailResponse, error) {
	postID, err := uuid.Parse(postIDRaw)

	if err != nil {
		ps.logger.Error("Error while parsing postID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	hasAccess, err := ps.storage.PostStore.HasAccessToPost(userID, postID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		ps.logger.Error("User has no access to this post")
		return nil, errors.InvalidPermissionError
	}

	post, err := ps.storage.PostStore.GetPostDetailsByID(postID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.logger.Error("Post not found")
			return nil, errors.NoPostsFoundError
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewPostDetailResponse(post)
	return result, nil
}
func (ps *PostService) CreatePost(userID uuid.UUID, dto dto.CreatePostDTO) error {

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

	err := ps.storage.PostStore.CreatePost(post)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ps *PostService) UpdatePost(userID uuid.UUID, postIDRaw string, dto dto.UpdatePostDTO) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		logger.Error("Error while parsing postID", zap.Error(err))
		return errors.InternalServerError.Wrap(err)
	}

	existPost, err := ps.storage.PostStore.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.logger.Error("Post not found")
			return errors.PostNotFoundError
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existPost.UserID != userID {
		ps.logger.Error("Request userid and post userid is different")
		return errors.InvalidPermissionError
	}
	post := &model.Post{
		ID:      postID,
		Content: dto.Content,
	}

	err = ps.storage.PostStore.UpdatePost(post)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ps *PostService) DeletePost(userID uuid.UUID, postIDRaw string) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		logger.Error("Error while parsing postID", zap.Error(err))
		return errors.InvalidIDFormatError
	}
	post, err := ps.storage.PostStore.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.logger.Error("Post not found")
			return errors.PostNotFoundError
		}
		return errors.InternalServerError.Wrap(err)
	}
	if post.UserID != userID {
		ps.logger.Error("Request userid and post userid is different")
		return errors.InvalidPermissionError
	}

	err = ps.storage.PostStore.DeletePost(postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
