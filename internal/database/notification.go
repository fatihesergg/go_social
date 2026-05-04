package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
)

type BaseNotificationStore interface {
	CreateNotification(ctx context.Context, notification model.Notification) error
	GetNotifications(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]model.Notification, error)
}

type NotificationStore struct {
	DB *sql.DB
}

func NewNotificationStore(db *sql.DB) BaseNotificationStore {
	return &NotificationStore{
		DB: db,
	}
}

func (ns *NotificationStore) CreateNotification(ctx context.Context, notification model.Notification) error {
	query := "INSERT INTO notifications(id,user_id,notification_type,message,payload,is_read) VALUES ($1,$2,$3,$4,$5,$6)"
	result, err := ns.DB.ExecContext(ctx, query, notification.ID, notification.UserID, notification.Notification_type.String(), notification.Message, notification.Payload, false)
	if err != nil {
		return fmt.Errorf("Error while creating notification: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error while getting affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("Error while creating notification,no rows affected")
	}

	return nil
}

func (ns *NotificationStore) GetNotifications(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]model.Notification, error) {
	query := "SELECT id,message,payload FROM notifications WHERE user_id = $1 ORDER BY  created_at DESC limit $2 offset $3"

	result, err := ns.DB.QueryContext(ctx, query, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, fmt.Errorf("Error while getting notifications: %w", err)
	}

	defer result.Close()

	notifications := []model.Notification{}
	for result.Next() {
		notification := model.Notification{}
		err := result.Scan(&notification.ID, &notification.Message, &notification.Payload)
		if err != nil {
			return nil, fmt.Errorf("Error while scanning notification result: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating result: %w", err)
	}
	return notifications, nil

}
