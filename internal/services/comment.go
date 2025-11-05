package services

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseCommentService interface {
	AddCommentPost(userID uuid.UUID, dto dto.CreateCommentDTO) error
	GetCommentsByPostID(userID uuid.UUID, postIDRaw string) ([]dto.CommentResponse, error)
	UpdateComment(userID uuid.UUID, commentIDRaw string, dto dto.UpdateCommentDTO) error
	DeleteComment(userID uuid.UUID, commentIDRaw string) error
}
type CommentService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewCommentService(storage *database.Storage, logger *zap.Logger) BaseCommentService {
	return &CommentService{storage: storage, logger: logger.Named("comment_service")}
}

func (cs *CommentService) AddCommentPost(userID uuid.UUID, dto dto.CreateCommentDTO) error {

	hasAccess, err := cs.storage.PostStore.HasAccessToPost(userID, dto.PostID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		cs.logger.Error("User has no access to post")
		return errors.InvalidPermissionError
	}

	comment := &model.Comment{
		ID:      uuid.New(),
		PostID:  dto.PostID,
		UserID:  userID,
		Content: dto.Content,
	}

	err = cs.storage.CommentStore.CreateComment(comment)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (cs *CommentService) GetCommentsByPostID(userID uuid.UUID, postIDRaw string) ([]dto.CommentResponse, error) {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		cs.logger.Error("Error while parsing uuid", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	hasAccess, err := cs.storage.PostStore.HasAccessToPost(userID, postID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		cs.logger.Error("User has no access to this post")
		return nil, errors.InvalidPermissionError
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
		cs.logger.Error("Error while parsing commentID", zap.Error(err))
		return errors.InvalidIDFormatError

	}

	hasAccess, err := cs.storage.CommentStore.HasAccessToComment(userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		cs.logger.Error("User has no access to this comment")
		return errors.InvalidPermissionError
	}

	comment, err := cs.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if comment == nil {
		return errors.CommentNotFoundError
	}

	if comment.UserID != userID {
		cs.logger.Error("Request userid and comment userid is different")
		return errors.InvalidPermissionError
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
		cs.logger.Error("Error while parsing commentID", zap.Error(err))
		return errors.InvalidIDFormatError
	}

	hasAccess, err := cs.storage.CommentStore.HasAccessToComment(userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		cs.logger.Error("User has no access to this comment")
		return errors.InvalidPermissionError
	}

	comment, err := cs.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if comment == nil {
		return errors.CommentNotFoundError
	}

	if comment.UserID != userID {
		cs.logger.Error("Request userid and comment userid is different")
		return errors.InvalidPermissionError
	}

	err = cs.storage.CommentStore.DeleteComment(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
