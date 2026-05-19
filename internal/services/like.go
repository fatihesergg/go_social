package services

import (
	"context"
	"fmt"

	"github.com/fatihesergg/go_social/internal/broker"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseLikeService interface {
	LikePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error
	LikeComment(ctx context.Context, userID uuid.UUID, commentRawID string) error
	LikeReply(ctx context.Context, userID uuid.UUID, replyRawID string) error
	UnlikePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error
	UnlikeComment(ctx context.Context, userID uuid.UUID, commentRawID string) error
	UnlikeReply(ctx context.Context, userID uuid.UUID, replyRawID string) error
}

type LikeService struct {
	likeStore    database.BaseLikeStore
	postStore    database.BasePostStore
	commentStore database.BaseCommentStore
	publisher    broker.EventPublisher
}

func NewLikeService(likeStore database.BaseLikeStore, postStore database.BasePostStore, commentStore database.BaseCommentStore, publisher broker.EventPublisher) BaseLikeService {
	return &LikeService{
		likeStore:    likeStore,
		postStore:    postStore,
		publisher:    publisher,
		commentStore: commentStore,
	}
}
func (ls *LikeService) LikePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error {
	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("error while parsing postID: %w", err))
	}

	hasAccess, err := ls.postStore.HasAccessToPost(ctx, userID, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this post"))
	}

	liked, err := ls.likeStore.IsPostLiked(ctx, postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if liked {
		return errors.AlreadyPostLikeError.Wrap(fmt.Errorf("user already liked this post"))
	}

	err = ls.likeStore.LikePost(ctx, &model.PostLike{
		ID:     uuid.New(),
		PostID: postID,
		UserID: userID,
	})

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	event := dto.PostLikedEvent{
		LikerID: userID,
		PostID:  postID,
	}
	err = ls.publisher.Publish(ctx, "post", "like_event", event)

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	return nil
}

func (ls *LikeService) UnlikePost(ctx context.Context, userID uuid.UUID, postIDRaw string) error {

	postID, err := uuid.Parse(postIDRaw)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing uuid: %w", err))
	}

	hasAccess, err := ls.postStore.HasAccessToPost(ctx, userID, postID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this post"))
	}

	liked, err := ls.likeStore.IsPostLiked(ctx, postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if !liked {
		return errors.PostNotLikedError.Wrap(fmt.Errorf("user has not liked post yet"))
	}

	err = ls.likeStore.UnlikePost(ctx, postID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}

func (ls *LikeService) LikeComment(ctx context.Context, userID uuid.UUID, commentRawID string) error {

	commentID, err := uuid.Parse(commentRawID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentID: %w", err))
	}
	hasAccess, err := ls.commentStore.HasAccessToComment(ctx, userID, commentID)

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	existLike, err := ls.likeStore.IsCommentLiked(ctx, commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if existLike {
		return errors.AlreadyCommentLikeError
	}
	err = ls.likeStore.LikeComment(ctx, &model.CommentLike{
		ID:        uuid.New(),
		CommentID: commentID,
		UserID:    userID,
	})
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ls *LikeService) UnlikeComment(ctx context.Context, userID uuid.UUID, commentRawID string) error {

	commentID, err := uuid.Parse(commentRawID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing commentID: %w", err))
	}
	hasAccess, err := ls.commentStore.HasAccessToComment(ctx, userID, commentID)

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !hasAccess {
		return errors.InvalidPermissionError.Wrap(fmt.Errorf("user has no access to this comment"))
	}

	existLike, err := ls.likeStore.IsCommentLiked(ctx, commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !existLike {
		return errors.CommentNotLikedError.Wrap(fmt.Errorf("user has not liked comment yet"))
	}

	err = ls.likeStore.UnlikeComment(ctx, commentID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}

func (ls *LikeService) LikeReply(ctx context.Context, userID uuid.UUID, replyRawID string) error {
	replyID, err := uuid.Parse(replyRawID)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("error while parsing replyID: %w", err))
	}

	liked, err := ls.likeStore.IsReplyLiked(ctx, replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	if liked {
		return errors.AlreadyReplyLikeError.Wrap(fmt.Errorf("user already liked reply"))
	}

	err = ls.likeStore.LikeReply(ctx, &model.ReplyLike{
		ID:      uuid.New(),
		ReplyID: replyID,
		UserID:  userID,
	})

	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}
	return nil
}
func (ls *LikeService) UnlikeReply(ctx context.Context, userID uuid.UUID, replyRawID string) error {
	replyID, err := uuid.Parse(replyRawID)
	if err != nil {
		return errors.InvalidIDFormatError.Wrap(fmt.Errorf("error while parsing replyID: %w", err))
	}

	existLike, err := ls.likeStore.IsReplyLiked(ctx, replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)
	}

	if !existLike {
		return errors.ReplyNotLikedError.Wrap(fmt.Errorf("user has not liked reply yet"))
	}

	err = ls.likeStore.UnlikeReply(ctx, replyID, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(err)

	}
	return nil
}
