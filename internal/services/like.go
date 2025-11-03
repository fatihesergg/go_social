package services

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseLikeService interface {
	LikePost(userID uuid.UUID, postIDRaw string) error
	LikeComment(userID uuid.UUID, commentRawID string) error
	LikeReply(userID uuid.UUID, replyRawID string) error
	UnlikePost(userID uuid.UUID, postIDRaw string) error
	UnlikeComment(userID uuid.UUID, commentRawID string) error
	UnlikeReply(userID uuid.UUID, replyRawID string) error
}

type LikeService struct {
	storage *database.Storage
}

func NewLikeService(storage *database.Storage) BaseLikeService {
	return &LikeService{storage: storage}
}
func (ls *LikeService) LikePost(userID uuid.UUID, postIDRaw string) error {
	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	hasAccess, err := ls.storage.PostStore.HasAccessToPost(userID, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError
	}

	liked, err := ls.storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if liked {
		return errors.AlreadyPostLikeError
	}

	err = ls.storage.LikeStore.LikePost(&model.PostLike{
		ID:     uuid.New(),
		PostID: postID,
		UserID: userID,
	})

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}

func (ls *LikeService) UnlikePost(userID uuid.UUID, postIDRaw string) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError
	}

	hasAccess, err := ls.storage.PostStore.HasAccessToPost(userID, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError
	}

	liked, err := ls.storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if !liked {
		return errors.PostNotLikedError
	}

	err = ls.storage.LikeStore.UnlikePost(postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}

func (ls *LikeService) LikeComment(userID uuid.UUID, commentRawID string) error {

	commentID, err := uuid.Parse(commentRawID)
	if err != nil {
		return errors.InvalidIDFormatError
	}
	hasAccess, err := ls.storage.CommentStore.HasAccessToComment(userID, commentID)

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError
	}

	existLike, err := ls.storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if existLike {
		return errors.AlreadyCommentLikeError
	}
	err = ls.storage.LikeStore.LikeComment(&model.CommentLike{
		ID:        uuid.New(),
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ls *LikeService) UnlikeComment(userID uuid.UUID, commentRawID string) error {

	commentID, err := uuid.Parse(commentRawID)
	if err != nil {
		return errors.InvalidIDFormatError
	}
	hasAccess, err := ls.storage.CommentStore.HasAccessToComment(userID, commentID)

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError
	}

	existLike, err := ls.storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !existLike {
		return errors.CommentNotLikedError
	}

	err = ls.storage.LikeStore.UnlikeComment(commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}

func (ls *LikeService) LikeReply(userID uuid.UUID, replyRawID string) error {
	replyID, err := uuid.Parse(replyRawID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	liked, err := ls.storage.LikeStore.IsReplyLiked(replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if liked {
		return errors.AlreadyReplyLikeError
	}

	err = ls.storage.LikeStore.LikeReply(&model.ReplyLike{
		ID:      uuid.New(),
		ReplyID: replyID,
		UserID:  userID,
	})

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ls *LikeService) UnlikeReply(userID uuid.UUID, replyRawID string) error {
	replyID, err := uuid.Parse(replyRawID)
	if err != nil {
		return errors.InvalidIDFormatError
	}

	existLike, err := ls.storage.LikeStore.IsReplyLiked(replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !existLike {
		return errors.ReplyNotLikedError
	}

	err = ls.storage.LikeStore.UnlikeReply(replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}
