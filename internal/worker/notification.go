package worker

import (
	"encoding/json"
	"fmt"

	"github.com/fatihesergg/go_social/internal/dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type NotificationWorker struct {
	Channel *amqp.Channel
	Logger  *zap.Logger
}

func NewNotificationWorker(ch *amqp.Channel, logger *zap.Logger) NotificationWorker {
	return NotificationWorker{Channel: ch, Logger: logger}
}

func (nw *NotificationWorker) Consume() error {
	msgs, err := nw.Channel.Consume("notification_queue", "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Error while consuming notification_queue: %w", err)
	}
	for msg := range msgs {
		postLikedEvent := dto.PostLikedEvent{}
		err = json.Unmarshal(msg.Body, &postLikedEvent)
		if err != nil {
			nw.Logger.Error("Error while parsing postLikedEvent", zap.Error(err))
		}
		nw.Logger.Info(
			"Notification:",
			zap.String("LikerID", postLikedEvent.LikerID.String()),
			zap.String("PostID", postLikedEvent.PostID.String()))

	}
	return nil
}
