package controller

import "github.com/fatihesergg/go_social/internal/services"

type FollowController struct {
	FollowService services.BaseFollowService
}

func NewFollowController(userService services.BaseUserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}
