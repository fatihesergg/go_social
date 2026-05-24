package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseCommentService interface {
	AddCommentPost(ctx context.Context, userID uuid.UUID, dto dto.CreateCommentDTO) error
	GetCommentsByPostID(ctx context.Context, userID uuid.UUID, postIDRaw string) ([]model.Comment, error)
	UpdateComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.UpdateCommentDTO) error
	DeleteComment(ctx context.Context, userID uuid.UUID, commentIDRaw string) error
}
type CommentService struct {
	postStore    database.BasePostStore
	commentStore database.BaseCommentStore
}

func NewCommentService(postStore database.BasePostStore, commentStore database.BaseCommentStore) BaseCommentService {
	return &CommentService{
		postStore:    postStore,
		commentStore: commentStore,
	}
}

func (cs *CommentService) AddCommentPost(ctx context.Context, userID uuid.UUID, dto dto.CreateCommentDTO) error {

	hasAccess, err := cs.postStore.HasAccessToPost(ctx, userID, dto.PostID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to post"))
	}

	comment := &model.Comment{
		ID:      uuid.New(),
		PostID:  dto.PostID,
		UserID:  userID,
		Content: dto.Content,
	}

	err = cs.commentStore.CreateComment(ctx, comment)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (cs *CommentService) GetCommentsByPostID(ctx context.Context, userID uuid.UUID, postIDRaw string) ([]model.Comment, error) {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing uuid: %w", err))
	}

	hasAccess, err := cs.postStore.HasAccessToPost(ctx, userID, postID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this post"))
	}

	comments, err := cs.commentStore.GetCommentsByPostID(ctx, postID, userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.CommentNotFoundError.Wrap(err)
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return comments, nil
}
func (cs *CommentService) UpdateComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.UpdateCommentDTO) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentid: %w", err))

	}

	hasAccess, err := cs.commentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	comment, err := cs.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.CommentNotFoundError.Wrap(err)
		}
		return errors.InternalServerError.Wrap(err)
	}

	if comment.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("request userid and comment userid is different"))
	}
	comment.Content = dto.Content

	err = cs.commentStore.UpdateComment(ctx, comment)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}
func (cs *CommentService) DeleteComment(ctx context.Context, userID uuid.UUID, commentIDRaw string) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentid: %w", err))
	}

	hasAccess, err := cs.commentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	comment, err := cs.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.CommentNotFoundError.Wrap(err)
		}
		return errors.InternalServerError.Wrap(err)
	}

	if comment.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("request userid and comment userid is different"))
	}

	err = cs.commentStore.DeleteComment(ctx, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
