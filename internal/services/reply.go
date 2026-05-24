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
	GetCommentReplies(ctx context.Context, userID uuid.UUID, commentIDRaw string) ([]model.Reply, error)
	GetRepliesByParentID(ctx context.Context, userID uuid.UUID, parentIDRaw string) ([]model.Reply, error)
	ReplyComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error
	ReplyAReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error
	UpdateReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error
	DeleteReply(ctx context.Context, userUD uuid.UUID, replyIDRaw string) error
}

type ReplyService struct {
	replyStore   database.BaseReplyStore
	commentStore database.BaseCommentStore
}

func NewReplyService(replyStore database.BaseReplyStore, commentStore database.BaseCommentStore) BaseReplyService {
	return &ReplyService{replyStore: replyStore, commentStore: commentStore}
}
func (rp *ReplyService) GetCommentReplies(ctx context.Context, userID uuid.UUID, commentIDRaw string) ([]model.Reply, error) {
	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentid: %w", err))
	}

	hasAccess, err := rp.commentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	existComment, err := rp.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.CommentNotFoundError.Wrap(err)
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	replies, err := rp.replyStore.GetRepliesByCommentID(ctx, existComment.ID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoRepliesFoundError.Wrap(err)
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return replies, nil
}
func (rp *ReplyService) ReplyComment(ctx context.Context, userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentid: %w", err))
	}

	hasAccess, err := rp.commentStore.HasAccessToComment(ctx, userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	comment, err := rp.commentStore.GetCommentByID(ctx, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.CommentNotFoundError.Wrap(err)
		}
		return errors.InternalServerError.Wrap(err)
	}

	reply := &model.Reply{
		ID:        uuid.New(),
		CommentID: comment.ID,
		UserID:    userID,
		Message:   dto.Message,
	}

	err = rp.replyStore.CreateCommentReply(ctx, reply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) UpdateReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("error while parsing replyid: %w", err))
	}

	hasAccess, err := rp.replyStore.HasAccessToReply(ctx, userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this reply"))
	}

	existReply, err := rp.replyStore.GetReplyByID(ctx, replyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(err)
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existReply.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("request userid and reply userid is different"))
	}

	existReply.Message = dto.Message

	err = rp.replyStore.UpdateReply(ctx, existReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil

}

func (rp *ReplyService) DeleteReply(ctx context.Context, userID uuid.UUID, replyIDRaw string) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing replyid: %w", err))
	}

	hasAccess, err := rp.replyStore.HasAccessToReply(ctx, userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this reply"))
	}

	existReply, err := rp.replyStore.GetReplyByID(ctx, replyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(err)
		}
		return errors.InternalServerError.Wrap(err)
	}

	if existReply.UserID != userID {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("request userid and reply id is different"))
	}

	err = rp.replyStore.DeleteReply(ctx, existReply.ID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) GetRepliesByParentID(ctx context.Context, userID uuid.UUID, parentIDRaw string) ([]model.Reply, error) {

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing parentid: %w", err))
	}

	hasAccess, err := rp.replyStore.HasAccessToReply(ctx, userID, parentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return nil, errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this reply"))
	}

	replies, err := rp.replyStore.GetRepliesByParentID(ctx, parentID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ReplyNotFoundError.Wrap(err)
		}
		return nil, errors.InternalServerError.Wrap(err)
	}

	return replies, nil
}

func (rp *ReplyService) ReplyAReply(ctx context.Context, userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error {
	parentID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing parentid: %w", err))
	}

	hasAccess, err := rp.replyStore.HasAccessToReply(ctx, userID, parentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this reply"))
	}

	reply, err := rp.replyStore.GetReplyByID(ctx, parentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ReplyNotFoundError.Wrap(err)
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

	err = rp.replyStore.CreateNestedReply(ctx, nestedReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
