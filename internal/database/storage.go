package database

type Storage struct {
	UserStore    BaseUserStore
	PostStore    BasePostStore
	CommentStore BaseCommentStore
}

func NewPostgresStorage(userStore BaseUserStore, postStore BasePostStore, commentStore BaseCommentStore) *Storage {
	return &Storage{
		UserStore:    userStore,
		PostStore:    postStore,
		CommentStore: commentStore,
	}
}
