package database

type Storage struct {
	UserStore BaseUserStore
	PostStore BasePostStore
}

func NewPostgresStorage(userStore BaseUserStore, postStore BasePostStore) *Storage {
	return &Storage{
		UserStore: userStore,
		PostStore: postStore,
	}
}
