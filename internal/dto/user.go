package dto

type CreateUserDTO struct {
	Name     string  `json:"name" binding:"required,alphanum,lte=50"`
	LastName string  `json:"last_name" binding:"required,alphanum,lte=50"`
	Email    string  `json:"email" binding:"required,email,lte=100"`
	Avatar   *string `json:"avatar"`
	Username string  `json:"username" binding:"required,alphanum,lte=50"`
	Password string  `json:"password" binding:"required,lte=20"`
}

type LoginUserDTO struct {
	Email    string `json:"email" binding:"required,email,lte=100"`
	Password string `json:"password" binding:"required,lte=20"`
}

type ResetUserPasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required,lte=20"`
	NewPassword string `json:"new_password" binding:"required,lte=20"`
}
