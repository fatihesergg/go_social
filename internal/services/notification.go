package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseNotificationService interface {
	SavePostLikedNotification(ctx context.Context, userID uuid.UUID, postID string, isRead bool) error
}

type NotificationService struct {
	storage *database.Storage
}

func NewNotificationService(storage *database.Storage) BaseNotificationService {
	return &NotificationService{
		storage: storage,
	}
}

func (ns *NotificationService) SavePostLikedNotification(ctx context.Context, userID uuid.UUID, postID string, isRead bool) error {
	user, err := ns.storage.UserStore.GetUserByID(ctx, userID)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while getting user information: %w", err))
	}

	message := fmt.Sprintf("%s liked your post", user.Username)

	payload := map[string]interface{}{
		"liker_id": userID.String(),
		"post_id":  postID,
	}

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while parsing post uuid: %w", err))
	}
	post, err := ns.storage.PostStore.GetPostByID(ctx, postUUID)

	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while getting PostLikedNotification post: %w", err))
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while marshal payload: %w", err))
	}

	model := model.Notification{
		ID:                uuid.New(),
		UserID:            post.UserID,
		Message:           message,
		Payload:           jsonPayload,
		Notification_type: model.PostLikedNotification,
		IsRead:            isRead,
	}

	err = ns.storage.NotificationStore.CreateNotification(ctx, model)
	if err != nil {
		return errors.InternalServerError.Wrap(fmt.Errorf("Error while creating post_liked notification: %w", err))
	}
	return nil

}
