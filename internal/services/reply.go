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

type BaseReplyService interface {
	GetCommentReplies(ctx context.Context, userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error)
	GetRepliesByParentID(ctx context.Context, userID uuid.UUID, parentIDRaw string) ([]dto.ReplyResponse, error)
	ReplyComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error
	ReplyAReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error
	UpdateReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error
	DeleteReply(ctx context.Context, userUD uuid.UUID, replyIDRaw string) error
}

type ReplyService struct {
	storage *database.Storage
}

func NewReplyService(storage *database.Storage) BaseReplyService {
	return &ReplyService{storage: storage}
}
func (rp *ReplyService) GetCommentReplies(ctx context.Context, userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error) {
	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing commentID: %w", err))
	}

	hasAccess, err := rp.storage.CommentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this comment"))
	}

	existComment, err := rp.storage.CommentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.CommentNotFoundError.Wrap(fmt.Errorf("Comment not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	replies, err := rp.storage.ReplyStore.GetRepliesByCommentID(ctx, existComment.ID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoRepliesFoundError.Wrap(fmt.Errorf("No replies found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewReplyResponse(replies)

	return result, nil
}
func (rp *ReplyService) ReplyComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing commentID: %w", err))
	}

	hasAccess, err := rp.storage.CommentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this comment"))
	}

	comment, err := rp.storage.CommentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.CommentNotFoundError.Wrap(fmt.Errorf("Comment not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	reply := &model.Reply{
		ID:        uuid.New(),
		CommentID: comment.ID,
		UserID:    userID,
		Message:   dto.Message,
	}

	err = rp.storage.ReplyStore.CreateCommentReply(ctx, reply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) UpdateReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while parsing replyID: %w", err)).Wrap(err)
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(ctx, userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this reply"))
	}

	existReply, err := rp.storage.ReplyStore.GetReplyByID(ctx, replyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(fmt.Errorf("Reply not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existReply.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("Request userid and reply userid is different"))
	}

	existReply.Message = dto.Message

	err = rp.storage.ReplyStore.UpdateReply(ctx, existReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil

}

func (rp *ReplyService) DeleteReply(ctx context.Context, userID uuid.UUID, replyIDRaw string) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing replyID: %w", err))
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(ctx, userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this reply"))
	}

	existReply, err := rp.storage.ReplyStore.GetReplyByID(ctx, replyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(fmt.Errorf("Reply not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existReply.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("Request userid and reply id is different"))
	}

	err = rp.storage.ReplyStore.DeleteReply(ctx, existReply.ID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) GetRepliesByParentID(ctx context.Context, userID uuid.UUID, parentIDRaw string) ([]dto.ReplyResponse, error) {

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing parentID: %w", err))
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(ctx, userID, parentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this reply"))
	}

	replies, err := rp.storage.ReplyStore.GetRepliesByParentID(ctx, parentID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ReplyNotFoundError.Wrap(fmt.Errorf("Reply not found"))
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	result := dto.NewReplyResponse(replies)
	return result, nil
}

func (rp *ReplyService) ReplyAReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error {
	parentID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("Error while parsing parentID: %w", err))
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(ctx, userID, parentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("User has no access to this reply"))
	}

	reply, err := rp.storage.ReplyStore.GetReplyByID(ctx, parentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(fmt.Errorf("Reply not found"))
		}
		return errors.InternalServerError.Wrap(err)
	}

	if reply == nil {
		return errors.ReplyNotFoundError
	}

	nestedReply := &model.Reply{
		ID:       uuid.New(),
		ParentID: parentID,
		UserID:   userID,
		Message:  dto.Message,
	}

	err = rp.storage.ReplyStore.CreateNestedReply(ctx, nestedReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
