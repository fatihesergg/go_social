package services

import (
	"context"
	"fmt"

	"github.com/fatihesergg/go_social/internal/appError"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseNotificationService interface {
	GetNotifications(ctx context.Context, userID uuid.UUID, pagination database.Pagination) ([]model.Notification, error)
}

type NotificationService struct {
	notificationStore database.BaseNotificationStore
	logger            *zap.Logger
}

func NewNotificationService(notificationStore database.BaseNotificationStore, logger *zap.Logger) BaseNotificationService {
	return &NotificationService{
		notificationStore: notificationStore,
		logger:            logger,
	}
}

func (ns *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, pagination database.Pagination) ([]model.Notification, error) {
	notifications, err := ns.notificationStore.GetNotifications(ctx, userID, pagination)
	if err != nil {
		ns.logger.Error("Error while getting notifications", zap.Error(err))
		return nil, appError.InternalServerError.Wrap(fmt.Errorf("error while getting notifications: %w", err))
	}

	return notifications, nil
}
