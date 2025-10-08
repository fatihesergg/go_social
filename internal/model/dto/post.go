package dto

type CreatePostDTO struct {
	Content string `json:"content" binding:"required,alphanum,lte=500"`
	Image   string `json:"image"`
}

type UpdatePostDTO struct {
	Content string `json:"content" binding:"required,alphanum,lte=500"`
	Image   string `json:"image"`
}
