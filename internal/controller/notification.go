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

// GetNotification godoc
//
//	@Summary		Get notifications
//	@Description	Retrieve a list of notifications
//	@Tags			Notification
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Success		200		{array}		util.SuccessResponse{result=[]dto.NotificationResponse}
//	@Failure		400		{object}	util.ErrorResponse
//	@Failure		500		{object}	util.ErrorResponse
//	@Router			/notifications [get]
//	@Security		Bearer
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
	result := dto.NewNotificationResponse(notifications)

	c.JSON(200, util.SuccessResponse{Result: result, Message: "Notification fetched successfully"})
}
