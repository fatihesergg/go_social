package services

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseReplyService interface {
	GetCommentReplies(userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error)
	GetRepliesByParentID(userID uuid.UUID, parentIDRaw string) ([]dto.ReplyResponse, error)
	ReplyComment(userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error
	ReplyAReply(userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error
	UpdateReply(userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error
	DeleteReply(userUD uuid.UUID, replyIDRaw string) error
}

type ReplyService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewReplyService(storage *database.Storage, logger *zap.Logger) BaseReplyService {
	return &ReplyService{storage: storage, logger: logger.Named("reply_service")}
}
func (rp *ReplyService) GetCommentReplies(userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error) {
	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		logger.Error("Error while parsing commentID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	hasAccess, err := rp.storage.CommentStore.HasAccessToComment(userID, commentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this comment")
		return nil, errors.InvalidPermissionError
	}

	existComment, err := rp.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if existComment == nil {
		return nil, errors.CommentNotFoundError
	}

	replies, err := rp.storage.ReplyStore.GetRepliesByCommentID(existComment.ID, userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if replies == nil {
		return nil, errors.NoRepliesFoundError
	}

	result := dto.NewReplyResponse(replies)

	return result, nil
}
func (rp *ReplyService) ReplyComment(userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error {

	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		logger.Error("Error while parsing commentID", zap.Error(err))
		return errors.InvalidIDFormatError
	}

	hasAccess, err := rp.storage.CommentStore.HasAccessToComment(userID, commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this comment")
		return errors.InvalidPermissionError
	}

	comment, err := rp.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if comment == nil {
		return errors.CommentNotFoundError
	}

	reply := &model.Reply{
		ID:        uuid.New(),
		CommentID: comment.ID,
		UserID:    userID,
		Message:   dto.Message,
	}

	err = rp.storage.ReplyStore.CreateCommentReply(reply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) UpdateReply(userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		logger.Error("Error while parsing replyID", zap.Error(err))
		return errors.InternalServerError.Wrap(err)
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this reply")
		return errors.InvalidPermissionError
	}

	existReply, err := rp.storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if existReply == nil {
		return errors.ReplyNotFoundError
	}

	if existReply.UserID != userID {
		rp.logger.Error("Request userid and reply userid is different")
		return errors.InvalidPermissionError
	}

	existReply.Message = dto.Message

	err = rp.storage.ReplyStore.UpdateReply(existReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil

}

func (rp *ReplyService) DeleteReply(userID uuid.UUID, replyIDRaw string) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		logger.Error("Error while parsing replyID", zap.Error(err))
		return errors.InvalidIDFormatError
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(userID, replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this reply")
		return errors.InvalidPermissionError
	}

	existReply, err := rp.storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if existReply == nil {
		return errors.ReplyNotFoundError
	}

	if existReply.UserID != userID {
		rp.logger.Error("Request userid and reply id is different")
		return errors.InvalidPermissionError
	}

	err = rp.storage.ReplyStore.DeleteReply(existReply.ID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) GetRepliesByParentID(userID uuid.UUID, parentIDRaw string) ([]dto.ReplyResponse, error) {

	parentID, err := uuid.Parse(parentIDRaw)
	if err != nil {
		logger.Error("Error while parsing parentID", zap.Error(err))
		return nil, errors.InvalidIDFormatError
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(userID, parentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this reply")
		return nil, errors.InvalidPermissionError
	}

	replies, err := rp.storage.ReplyStore.GetRepliesByParentID(parentID, userID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}
	if replies == nil {
		return nil, errors.NoRepliesFoundError
	}
	result := dto.NewReplyResponse(replies)
	return result, nil
}

func (rp *ReplyService) ReplyAReply(userID uuid.UUID, replyIDRaw string, dto dto.CreateReply) error {
	parentID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		logger.Error("Error while parsing parentID", zap.Error(err))
		return errors.InvalidIDFormatError
	}

	hasAccess, err := rp.storage.ReplyStore.HasAccessToReply(userID, parentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		rp.logger.Error("User has no access to this reply")
		return errors.InvalidPermissionError
	}

	reply, err := rp.storage.ReplyStore.GetReplyByID(parentID)
	if err != nil {
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

	err = rp.storage.ReplyStore.CreateNestedReply(nestedReply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
