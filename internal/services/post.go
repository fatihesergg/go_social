package services

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
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
}

func NewPostService(storage *database.Storage) BasePostService {
	return &PostService{storage: storage}
}
func (ps *PostService) GetAllPosts(userID uuid.UUID, limit, offset, query string) ([]dto.AllPostResponse, error) {
	pagination := database.NewPagination(limit, offset)
	search := database.NewSearch(query)

	posts, err := ps.storage.PostStore.GetPosts(pagination, search, userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}
	if posts == nil {
		return nil, errors.NoPostsFoundError
	}
	result := dto.NewAllPostResponse(posts)
	return result, nil

}

func (ps *PostService) GetPostByID(userID uuid.UUID, postIDRaw string) (*dto.PostDetailResponse, error) {
	postID, err := uuid.Parse(postIDRaw)

	if err != nil {
		return nil, errors.InvalidIDFormatError
	}

	post, err := ps.storage.PostStore.GetPostDetailsByID(postID, userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}
	if post == nil {
		return nil, errors.PostNotFoundError
	}
	result := dto.NewPostDetailResponse(post)
	return result, nil
}
func (ps *PostService) CreatePost(userID uuid.UUID, dto dto.CreatePostDTO) error {

	post := &model.Post{
		Content: dto.Content,
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
		return errors.InternalServerError.Wrap(err)
	}

	existPost, err := ps.storage.PostStore.GetPostByID(postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if existPost == nil {
		return errors.PostNotFoundError
	}
	if existPost.UserID != userID {
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
		return errors.InvalidIDFormatError
	}
	post, err := ps.storage.PostStore.GetPostByID(postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if post == nil {
		return errors.PostNotFoundError
	}
	if post.UserID != userID {
		return errors.InvalidPermissionError
	}

	err = ps.storage.PostStore.DeletePost(postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
