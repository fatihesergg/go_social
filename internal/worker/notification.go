package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/ws"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type NotificationWorker struct {
	Channel           *amqp.Channel
	Logger            *zap.Logger
	hub               *ws.WsHub
	postStore         database.BasePostStore
	userStore         database.BaseUserStore
	commentStore      database.BaseCommentStore
	replyStore        database.BaseReplyStore
	notificationStore database.BaseNotificationStore
}

func NewNotificationWorker(ch *amqp.Channel,
	logger *zap.Logger,
	hub *ws.WsHub,
	postStore database.BasePostStore,
	userStore database.BaseUserStore,
	commentStore database.BaseCommentStore,
	replyStore database.BaseReplyStore,
	notificationStore database.BaseNotificationStore) NotificationWorker {
	return NotificationWorker{Channel: ch,
		Logger:            logger,
		hub:               hub,
		postStore:         postStore,
		userStore:         userStore,
		notificationStore: notificationStore,
		commentStore:      commentStore,
		replyStore:        replyStore,
	}
}

func (nw *NotificationWorker) ConsumePostLike(ctx context.Context) error {
	msgs, err := nw.Channel.Consume("post_like_queue", "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error while consuming post_like_queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			nw.handlePostLike(ctx, msg.Body)

		}
	}
}

func (nw *NotificationWorker) handlePostLike(ctx context.Context, event []byte) {
	data := dto.PostLikedEvent{}
	err := json.Unmarshal(event, &data)

	if err != nil {
		nw.Logger.Error("error while parsing postLikedEvent", zap.Error(err))
		return
	}

	nw.Logger.Info(
		"post like event",
		zap.String("likerID", data.LikerID.String()),
		zap.String("postID", data.PostID.String()),
	)

	dbContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	post, err := nw.postStore.GetPostByID(dbContext, data.PostID)
	if err != nil {
		nw.Logger.Error("Error while getting post information: %w", zap.Error(err))
		return
	}

	liker, err := nw.userStore.GetUserByID(dbContext, data.LikerID)
	if err != nil {
		nw.Logger.Error("Error while getting user information: %w", zap.Error(err))
		return

	}

	payload := model.PostLikePayload{
		PostID:     post.ID.String(),
		PostUserID: post.UserID.String(),
		LikerID:    liker.ID.String(),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		nw.Logger.Error("error while parsing json payload", zap.String("worker", "post.like"), zap.Error(err))
		return
	}

	notification := model.Notification{
		ID:                uuid.New(),
		UserID:            post.UserID,
		Message:           fmt.Sprintf("%s liked your post", liker.Username),
		Payload:           jsonPayload,
		IsRead:            false,
		Notification_type: model.PostLikeNotificationType,
		Timestamp:         time.Now(),
	}

	err = nw.notificationStore.CreateNotification(dbContext, notification)
	if err != nil {
		nw.Logger.Error("error while saving notification to database", zap.Error(err))
		return
	}

	postNotify := model.PostLikeNotification{
		ID:        notification.ID,
		UserID:    post.UserID,
		Message:   fmt.Sprintf("%s liked your post", liker.Username),
		Payload:   payload,
		IsRead:    false,
		Timestamp: notification.Timestamp,
	}

	jsonNotify, err := json.Marshal(postNotify)
	if err != nil {
		nw.Logger.Error("error while marshalling post like response", zap.Error(err))
		return
	}

	nw.hub.Messages <- ws.WsData{
		Data:   jsonNotify,
		UserID: post.UserID,
	}
}

func (nw *NotificationWorker) ConsumeCommentLike(ctx context.Context) error {

	ch, err := nw.Channel.Consume("comment_like_queue", "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error while consuming comment_like_queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			nw.handleCommentLike(ctx, msg.Body)
		}
	}
}

func (nw *NotificationWorker) handleCommentLike(ctx context.Context, event []byte) {
	data := dto.CommentLikedEvent{}

	err := json.Unmarshal(event, &data)
	if err != nil {
		nw.Logger.Error("error while unmarshalling comment like event", zap.Error(err))
		return
	}

	dbContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	liker, err := nw.userStore.GetUserByID(dbContext, data.LikerID)
	if err != nil {
		nw.Logger.Error("error while getting user", zap.Error(err))
		return
	}

	comment, err := nw.commentStore.GetCommentByID(dbContext, data.CommentID)
	if err != nil {
		nw.Logger.Error("error while getting comment", zap.Error(err))
		return
	}

	commentUser, err := nw.userStore.GetUserByID(dbContext, comment.UserID)
	if err != nil {
		nw.Logger.Error("error while getting comment user", zap.Error(err))
		return
	}

	post, err := nw.postStore.GetPostByID(dbContext, comment.PostID)
	if err != nil {
		nw.Logger.Error("error while getting comment post", zap.Error(err))
		return
	}

	payload := model.CommentLikePayload{
		LikerID:       liker.ID,
		PostID:        post.ID,
		CommentID:     comment.ID,
		CommentUserID: comment.UserID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		nw.Logger.Error("error while marshalling json payload", zap.Error(err))
		return
	}

	notification := model.Notification{
		ID:                uuid.New(),
		UserID:            commentUser.ID,
		Notification_type: model.CommentLikeNotificationType,
		Payload:           jsonPayload,
		Message:           fmt.Sprintf("%s liked your comment", liker.Username),
		IsRead:            false,
		Timestamp:         time.Now(),
	}

	err = nw.notificationStore.CreateNotification(dbContext, notification)
	if err != nil {
		nw.Logger.Error("error while creating notification", zap.Error(err))
		return
	}

	commentNotify := model.CommentLikeNotification{
		ID:                notification.ID,
		UserID:            comment.UserID,
		Notification_type: model.CommentLikeNotificationType,
		Payload:           payload,
		Message:           fmt.Sprintf("%s liked your comment", liker.Username),
		Timestamp:         notification.Timestamp,
	}

	jsonNotify, err := json.Marshal(commentNotify)
	if err != nil {
		nw.Logger.Error("error while marshalling comment like response", zap.Error(err))
		return
	}

	nw.hub.Messages <- ws.WsData{
		Data:   jsonNotify,
		UserID: comment.UserID,
	}

	commentNotify.Message = fmt.Sprintf("%s like %s's comment in your post", liker.Username, comment.User.Username)

	jsonNotify, err = json.Marshal(commentNotify)
	if err != nil {
		nw.Logger.Error("error while marshalling other comment like response", zap.Error(err))
		return
	}

	nw.hub.Messages <- ws.WsData{
		Data:   jsonNotify,
		UserID: post.UserID,
	}

}

func (nw *NotificationWorker) ConsumeReplyLike(ctx context.Context) error {

	ch, err := nw.Channel.Consume("reply_like_queue", "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error while consuming reply_like_queue: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			nw.handleReplyLike(ctx, msg.Body)
		}
	}
}

func (nw *NotificationWorker) handleReplyLike(ctx context.Context, event []byte) {

	data := dto.ReplyLikedEvent{}

	err := json.Unmarshal(event, &data)
	if err != nil {
		nw.Logger.Error("error while unmarshalling reply like event", zap.Error(err))
		return
	}

	dbContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	liker, err := nw.userStore.GetUserByID(dbContext, data.LikerID)
	if err != nil {
		nw.Logger.Error("error while getting liker user", zap.Error(err))
		return
	}

	reply, err := nw.replyStore.GetReplyByID(dbContext, data.ReplyID)
	if err != nil {
		nw.Logger.Error("error while getting reply", zap.Error(err))
		return
	}

	payload := model.ReplyLikePayload{
		ReplyID: data.ReplyID,
		LikerID: data.LikerID,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		nw.Logger.Error("error while marshalling reply like payload", zap.Error(err))
		return
	}

	notification := model.Notification{
		ID:                uuid.New(),
		UserID:            reply.UserID,
		Notification_type: model.ReplyLikeNotificationType,
		Message:           fmt.Sprintf("%s liked your reply", liker.Username),
		Timestamp:         time.Now(),
		Payload:           jsonPayload,
		IsRead:            false,
	}

	err = nw.notificationStore.CreateNotification(dbContext, notification)
	if err != nil {
		nw.Logger.Error("error while creating reply like notification", zap.Error(err))
		return
	}

	replyNotify := model.ReplyLikeNotification{
		ID:               notification.ID,
		UserID:           reply.UserID,
		NotificationType: model.ReplyLikeNotificationType,
		Message:          notification.Message,
		Payload:          payload,
		IsRead:           false,
		Timestamp:        notification.Timestamp,
	}

	jsonNotify, err := json.Marshal(replyNotify)
	if err != nil {
		nw.Logger.Error("error while marshalling reply like notification", zap.Error(err))
		return
	}

	nw.hub.Messages <- ws.WsData{
		UserID: reply.UserID,
		Data:   jsonNotify,
	}

}
