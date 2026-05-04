package services

import (
	"context"
	"fmt"

	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/errors"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BaseNotificationService interface {
	GetNotifications(ctx context.Context, userID uuid.UUID, pagination database.Pagination) ([]model.Notification, error)
}

type NotificationService struct {
	storage *database.Storage
	logger  *zap.Logger
}

func NewNotificationService(storage *database.Storage, logger *zap.Logger) BaseNotificationService {
	return &NotificationService{
		storage: storage,
		logger:  logger,
	}
}

func (ns *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, pagination database.Pagination) ([]model.Notification, error) {
	notifications, err := ns.storage.NotificationStore.GetNotifications(ctx, userID, pagination)
	if err != nil {
		ns.logger.Error("Error while getting notifications", zap.Error(err))
		return nil, errors.InternalServerError.Wrap(fmt.Errorf("Error while getting notifications: %w", err))
	}

	return notifications, nil
}
