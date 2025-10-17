package services

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseReplyService interface {
	GetCommentReplies(userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error)
	ReplyComment(userID uuid.UUID, commentIDRaw string, dto dto.CreateReply) error
	UpdateReply(userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error
	DeleteReply(userUD uuid.UUID, replyIDRaw string) error
}

type ReplyService struct {
	storage *database.Storage
}

func NewReplyService(storage *database.Storage) BaseReplyService {
	return &ReplyService{storage: storage}
}
func (rp *ReplyService) GetCommentReplies(userID uuid.UUID, commentIDRaw string) ([]dto.ReplyResponse, error) {
	commentID, err := uuid.Parse(commentIDRaw)
	if err != nil {
		return nil, errors.InvalidIDFormatError
	}

	existComment, err := rp.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return nil, errors.InternalServerError.Wrap(err)
	}

	if existComment == nil {
		return nil, errors.CommentNotFoundError
	}

	replies, err := rp.storage.ReplyStore.GetRepliesByCommentID(existComment.ID)
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
		return errors.InvalidIDFormatError
	}

	comment, err := rp.storage.CommentStore.GetCommentByID(commentID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if comment == nil {
		return errors.CommentNotFoundError
	}

	reply := &model.Reply{
		CommentID: comment.ID,
		UserID:    userID,
		Message:   dto.Message,
	}

	err = rp.storage.ReplyStore.CreateReply(reply)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (rp *ReplyService) UpdateReply(userID uuid.UUID, replyIDRaw string, dto dto.UpdateReply) error {

	replyID, err := uuid.Parse(replyIDRaw)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	existReply, err := rp.storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if existReply == nil {
		return errors.ReplyNotFoundError
	}

	if existReply.UserID != userID {
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
		return errors.InvalidIDFormatError
	}

	existReply, err := rp.storage.ReplyStore.GetReplyByID(replyID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if existReply == nil {
		return errors.ReplyNotFoundError
	}

	if existReply.UserID != userID {
		return errors.InvalidPermissionError
	}

	err = rp.storage.ReplyStore.DeleteReply(existReply.ID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
