package database

type Storage struct {
	UserStore    BaseUserStore
	PostStore    BasePostStore
	CommentStore BaseCommentStore
	FollowStore  BaseFollowStore
}

func NewPostgresStorage(userStore BaseUserStore, postStore BasePostStore, commentStore BaseCommentStore, followStore BaseFollowStore) *Storage {
	return &Storage{
		UserStore:    userStore,
		PostStore:    postStore,
		CommentStore: commentStore,
		FollowStore:  followStore,
	}
}
