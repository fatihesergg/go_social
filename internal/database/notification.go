package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fatihesergg/go_social/internal/model"
)

type BaseNotificationStore interface {
	CreateNotification(ctx context.Context, notification model.Notification) error
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
