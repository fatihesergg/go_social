package controller

import (
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationController struct {
	notificationService services.BaseNotificationService
}

func NewNotificationController(notificationService services.BaseNotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

func (nc *NotificationController) GetNotification(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	limit := c.Query("limit")
	offset := c.Query("offset")

	pagination := database.NewPagination(limit, offset)
	notifications, err := nc.notificationService.GetNotifications(c.Request.Context(), userID, pagination)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	data := dto.NewNotificationResponse(notifications)

	c.JSON(200, util.SuccessResultResponse{Result: data, Message: "Notification fetched successfully"})
}
