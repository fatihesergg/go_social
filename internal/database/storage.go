package database

type Storage struct {
	UserStore BaseUserStore
}

func NewPostgresStorage(userStore BaseUserStore) *Storage {
	return &Storage{
		UserStore: userStore,
	}
}
