package services

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentService interface {
	AddCommentPost(userID uuid.UUID, dto dto.CreateCommentDTO) error
	GetCommentsByPostID(userID uuid.UUID, postIDRaw string) ([]dto.CommentResponse, error)
	UpdateComment(userID uuid.UUID, commentIDRaw string, dto dto.UpdateCommentDTO) error
	DeleteComment(userID uuid.UUID, commentIDRaw string) error
}
type CommentService struct {
	storage *database.Storage
}

func NewCommentService(storage *database.Storage) BaseCommentService {
	return &CommentService{storage: storage}
}

func (cs *CommentService) AddCommentPost(userID uuid.UUID, dto dto.CreateCommentDTO) error {

	comment := &model.Comment{
		ID:      uuid.New(),
		PostID:  dto.PostID,
		UserID:  userID,
		Content: dto.Content,
	}

	err := cs.storage.CommentStore.CreateComment(comment)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (cs *CommentService) GetCommentsByPostID(userID uuid.UUID, postIDRaw string) ([]dto.CommentResponse, error) {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError
	}
	comments, err := cs.storage.CommentStore.GetCommentsByPostID(postID, userID)

	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if comments == nil {
		return nil, errors.CommentNotFoundError
	}
	result := dto.NewCommentResponse(comments)
	return result, nil
}
func (cs *CommentService) UpdateComment(userID uuid.UUID, commentIDRaw string, dto dto.UpdateCommentDTO) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError
	}

	comment, err := cs.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if comment == nil {
		return errors.CommentNotFoundError
	}
	comment.Content = dto.Content

	err = cs.storage.CommentStore.UpdateComment(comment)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}
func (cs *CommentService) DeleteComment(userID uuid.UUID, commentIDRaw string) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError
	}

	comment, err := cs.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if comment == nil {
		return errors.CommentNotFoundError
	}

	if comment.UserID != userID {
		return errors.InvalidPermissionError
	}

	err = cs.storage.CommentStore.DeleteComment(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
