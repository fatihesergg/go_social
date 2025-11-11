package database

type Storage struct {
	UserStore    BaseUserStore
	PostStore    BasePostStore
	CommentStore BaseCommentStore
	FollowStore  BaseFollowStore
	FeedStore    BaseFeedStore
	LikeStore    BaseLikeStore
	ReplyStore   BaseReplyStore
	TagStore     BaseTagStore
}

func NewPostgresStorage(userStore BaseUserStore, postStore BasePostStore, commentStore BaseCommentStore, followStore BaseFollowStore, feedStore BaseFeedStore, likeStore BaseLikeStore, replyStore BaseReplyStore, tagStore BaseTagStore) *Storage {
	return &Storage{
		UserStore:    userStore,
		PostStore:    postStore,
		CommentStore: commentStore,
		FollowStore:  followStore,
		FeedStore:    feedStore,
		LikeStore:    likeStore,
		ReplyStore:   replyStore,
		TagStore:     tagStore,
	}
}
