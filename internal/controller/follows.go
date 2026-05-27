package controller

import (
	"github.com/fatihesergg/go_social/internal/appError"
	"github.com/fatihesergg/go_social/internal/dto"
	"github.com/fatihesergg/go_social/internal/model"
	_ "github.com/fatihesergg/go_social/internal/model"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowController struct {
	followService services.BaseFollowService
}

func NewFollowController(followService services.BaseFollowService) *FollowController {
	return &FollowController{
		followService: followService,
	}
}

// AcceptFollowRequest godoc
//
//	@Summary		Accept follow request
//	@Description	Accept follow request
//	@Tags			Follows
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		201	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/accept [post]
func (fc *FollowController) AcceptFollowRequest(c *gin.Context) {
	senderID := c.Param("id")

	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	requesterID := c.MustGet("userID").(uuid.UUID)

	requestDto := dto.RespondFollowRequest{
		SenderID:    senderUUID,
		ResponderID: requesterID,
	}

	err = fc.followService.AcceptFollowRequest(c.Request.Context(), requestDto)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessResponse{Message: "follow request accepted successfully"})

}

// RejectFollowRequest godoc
//
//	@Summary		reject follow request
//	@Description	reject follow request
//	@Tags			Follows
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		201	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/reject [post]
func (fc *FollowController) RejectFollowRequest(c *gin.Context) {

	senderID := c.Param("id")

	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)

	requestDto := dto.RespondFollowRequest{
		SenderID:    senderUUID,
		ResponderID: userID,
	}

	err = fc.followService.RejectFollowRequest(c.Request.Context(), requestDto)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessResponse{Message: "follow request rejected successfully"})

}

// GetFollowerByUserID godoc
//
//	@Summary		Get followers of a user by user ID
//	@Description	Retrieve a list of followers for a specific user
//	@Tags			Follows
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	util.SuccessResponse{result=[]model.Follow}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/followers [get]
func (fc *FollowController) GetFollowerByUserID(c *gin.Context) {
	id := c.Param("id")

	followUUID, err := uuid.Parse(id)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	meID := c.MustGet("userID").(uuid.UUID)

	followModel := model.Follow{
		UserID:   meID,
		FollowID: followUUID,
	}

	followers, err := fc.followService.GetFollowerByUserID(c.Request.Context(), followModel)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(200, util.SuccessResponse{Message: "followers fetched successfully", Result: followers})
}

// GetFollowingByUserID godoc
//
//	@Summary		Get followings of a user by user ID
//	@Description	Retrieve a list of users that a specific user is following
//	@Tags			Follows
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	util.SuccessResponse{result=[]model.Follow}
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		404	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/following [get]
func (fc *FollowController) GetFollowingByUserID(c *gin.Context) {
	id := c.Param("id")

	followUUID, err := uuid.Parse(id)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	meID := c.MustGet("userID").(uuid.UUID)

	followModel := model.Follow{
		UserID:   meID,
		FollowID: followUUID,
	}
	followings, err := fc.followService.GetFollowingByUserID(c.Request.Context(), followModel)

	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(200, util.SuccessResponse{Message: "user following fetched successfully", Result: followings})
}

// FollowUser godoc
//
//	@Summary		Follow request
//	@Description	Send follow request to user with user ID
//	@Tags			Follows
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	true	"User ID to send follow request"
//	@Success		204	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/follow [post]
func (fc *FollowController) SendFollowRequest(c *gin.Context) {
	followID := c.Param("id")

	followUUID, err := uuid.Parse(followID)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	meID := c.MustGet("userID").(uuid.UUID)

	requestDto := dto.SendFollowRequest{
		RequesterID: meID,
		FollowID:    followUUID,
	}

	err = fc.followService.SendFollowRequest(c.Request.Context(), requestDto)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}

	c.JSON(201, util.SuccessResponse{Message: "follow request send successfully"})
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			Follows
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	true	"User ID to unfollow"
//	@Success		200	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/unfollow [post]
func (fc *FollowController) UnfollowUser(c *gin.Context) {
	id := c.Param("id")

	followID, err := uuid.Parse(id)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	meID := c.MustGet("userID").(uuid.UUID)

	requestDto := dto.UnfollowRequest{
		RequesterID: meID,
		FollowID:    followID,
	}

	err = fc.followService.UnFollowUser(c.Request.Context(), requestDto)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResponse{Message: "unfollowed successfully"})
}

// CancelFollowRequest godoc
//
//	@Summary		Cancel follow request
//	@Description	Cancel follow request with UserID
//	@Tags			Follows
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	true	"User ID to cancel follow request"
//	@Success		200	{object}	util.SuccessResponse
//	@Failure		400	{object}	util.ErrorResponse
//	@Failure		500	{object}	util.ErrorResponse
//	@Security		Bearer
//	@Router			/follows/{id}/cancel [post]
func (fc *FollowController) CancelFollowRequest(c *gin.Context) {

	id := c.Param("id")

	followID, err := uuid.Parse(id)
	if err != nil {
		util.WriteAppError(c, appError.InvalidIDFormatError)
		return
	}

	meID := c.MustGet("userID").(uuid.UUID)

	requestDto := dto.CancelFollowRequest{
		RequesterID: meID,
		FollowID:    followID,
	}

	err = fc.followService.CancelFollowRequest(c.Request.Context(), requestDto)
	if err != nil {
		util.WriteAppError(c, err)
		return
	}
	c.JSON(200, util.SuccessResponse{Message: "follow request cancelled"})

}
