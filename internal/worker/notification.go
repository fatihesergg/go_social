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
	notificationStore database.BaseNotificationStore
}

func NewNotificationWorker(ch *amqp.Channel,
	logger *zap.Logger,
	hub *ws.WsHub,
	postStore database.BasePostStore,
	userStore database.BaseUserStore,
	notificationStore database.BaseNotificationStore) NotificationWorker {
	return NotificationWorker{Channel: ch,
		Logger:            logger,
		hub:               hub,
		postStore:         postStore,
		userStore:         userStore,
		notificationStore: notificationStore,
	}
}

func (nw *NotificationWorker) Consume(ctx context.Context) error {
	msgs, err := nw.Channel.Consume("post_liked", "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error while consuming notification_queue: %w", err)
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
		nw.Logger.Error("Error while parsing postLikedEvent", zap.Error(err))
		return
	}

	nw.Logger.Info(
		"Notification:",
		zap.String("LikerID", data.LikerID.String()),
		zap.String("PostID", data.PostID.String()))

	dbContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	post, err := nw.postStore.GetPostByID(dbContext, data.PostID)
	if err != nil {
		nw.Logger.Error("Error while getting post information: %w", zap.Error(err))
		return
	}

	user, err := nw.userStore.GetUserByID(dbContext, data.LikerID)
	if err != nil {
		nw.Logger.Error("Error while getting user information: %w", zap.Error(err))
		return

	}

	postLikedNotify := model.Notification{
		ID:      uuid.New(),
		UserID:  post.UserID,
		Message: fmt.Sprintf("%s liked your post", user.Username),
		//NOTE: Send msg.Body as a payload for now.
		Payload: event,
		IsRead:  false,
	}

	err = nw.notificationStore.CreateNotification(dbContext, postLikedNotify)
	if err != nil {
		nw.Logger.Error("Error while saving notification to database", zap.Error(err))
		return
	}

	nw.hub.Broadcast <- ws.WsData{
		Data:   event,
		UserID: post.UserID,
	}
}
