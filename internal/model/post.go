package model 

type Post struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	

}