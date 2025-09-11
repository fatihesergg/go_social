package database

type Storage struct {
	UserStore    BaseUserStore
	PostStore    BasePostStore
	CommentStore BaseCommentStore
	FollowStore  BaseFollowStore
	FeedStore    BaseFeedStore
}

func NewPostgresStorage(userStore BaseUserStore, postStore BasePostStore, commentStore BaseCommentStore, followStore BaseFollowStore, feedStore BaseFeedStore) *Storage {
	return &Storage{
		UserStore:    userStore,
		PostStore:    postStore,
		CommentStore: commentStore,
		FollowStore:  followStore,
		FeedStore:    feedStore,
	}
}
