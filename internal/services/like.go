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
	UnlikePost(userID uuid.UUID, postIDRaw string) error
	UnlikeComment(userID uuid.UUID, commentRawID string) error
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

	liked, err := ls.storage.LikeStore.IsPostLiked(postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if liked {
		return errors.AlreadyPostLikeError
	}

	err = ls.storage.LikeStore.LikePost(&model.PostLike{
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

	existLike, err := ls.storage.LikeStore.IsCommentLiked(commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if existLike {
		return errors.AlreadyCommentLikeError
	}
	err = ls.storage.LikeStore.LikeComment(&model.CommentLike{
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
