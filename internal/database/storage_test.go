package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/fatihesergg/go_social/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB
var testStorage *Storage

func NewPostgresTestStorage() *Storage {
	godotenv.Load("../../.env")

	testDSN := os.Getenv("TEST_DB_URL")
	if testDSN == "" {
		panic("TEST_DB_URL environment variable is not set")
	}
	db, err := sql.Open("postgres", testDSN)
	if err != nil {
		fmt.Println("Failed to test database connection:", err.Error())
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Failed to create postgres driver:", err.Error())
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migration", "postgres", driver)
	if err != nil {
		fmt.Println("Failed to create migrate instance:", err.Error())
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		fmt.Println("Failed to run migrations:", err.Error())
		panic(err)
	}

	return &Storage{
		UserStore:    NewUserStore(db),
		PostStore:    NewPostStore(db),
		CommentStore: NewCommentStore(db),
		FollowStore:  NewFollowStore(db),
		FeedStore:    NewFeedStore(db),
		LikeStore:    NewLikeStore(db),
		ReplyStore:   NewReplyStore(db),
		TagStore:     NewTagStore(db),
	}
}

func cleanupAllTables() {
	tables := []string{"posts", "post_likes", "comments", "comment_likes", "users"}
	for _, table := range tables {
		if _, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			fmt.Printf("Error truncate table %s, %s \n", table, err.Error())
		}
	}
}

func createTestPagination(t *testing.T) Pagination {
	t.Helper()
	return Pagination{
		Limit:  20,
		Offset: 0,
	}
}

func createTestSearch(t *testing.T, query string) Search {
	t.Helper()
	return Search{
		Query: query,
	}
}

func createTestUser(t *testing.T, name, lastName, username, email, password string) *model.User {
	t.Helper()
	return &model.User{
		ID:       uuid.New(),
		Name:     name,
		LastName: lastName,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func createTestPost(t *testing.T, content string, userID uuid.UUID, visibility string) *model.Post {
	t.Helper()
	return &model.Post{
		ID:         uuid.New(),
		Content:    content,
		UserID:     userID,
		Visibility: visibility,
	}
}

func createTestComment(t *testing.T, content string, postID, userID uuid.UUID) *model.Comment {
	t.Helper()

	return &model.Comment{
		ID:      uuid.New(),
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
}

func createTestLikePost(t *testing.T, postID, userID uuid.UUID) *model.PostLike {
	t.Helper()
	return &model.PostLike{
		ID:     uuid.New(),
		PostID: postID,
		UserID: userID,
	}
}
func createTestLikeComment(t *testing.T, commentID, userID uuid.UUID) *model.CommentLike {
	t.Helper()
	return &model.CommentLike{
		ID:        uuid.New(),
		CommentID: commentID,
		UserID:    userID,
	}
}

func createTestLikeReply(t *testing.T, replyID uuid.UUID, userID uuid.UUID) *model.ReplyLike {
	t.Helper()
	return &model.ReplyLike{
		ReplyID: replyID,
		UserID:  userID,
	}
}

func createTestCommentReply(t *testing.T, commentID, userID uuid.UUID, message string) *model.Reply {
	t.Helper()
	return &model.Reply{
		ID:        uuid.New(),
		CommentID: commentID,
		UserID:    userID,
		Message:   message,
	}
}

func createTestNestedReply(t *testing.T, parentID, userID uuid.UUID, message string) *model.Reply {
	t.Helper()
	return &model.Reply{
		ID:       uuid.New(),
		ParentID: parentID,
		UserID:   userID,
		Message:  message,
	}
}

func createTestFollow(t *testing.T, userID, followID uuid.UUID) *model.Follow {
	t.Helper()
	return &model.Follow{
		ID:       uuid.New(),
		UserID:   userID,
		FollowID: followID,
	}
}

func createTestTags(tags ...string) []model.Tag {
	tagSlice := make([]model.Tag, len(tags))
	for i := 0; i < len(tags); i++ {
		tagSlice[i] = model.Tag{
			ID:   uuid.New(),
			Name: tags[i],
		}
	}
	return tagSlice
}

func TestUserStore_CreateUser(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existUser)
	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, existUser.ID)
	})
}

func TestUserStore_UpdateUser(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	user.Name = "test_update"
	user.LastName = "test_update"
	user.Username = "test_update"
	user.Email = "test_update@test.com"
	user.Password = "test_update"

	err = testStorage.UserStore.UpdateUser(ctx, user)
	assert.NoError(t, err)

	updatedUser, err := testStorage.UserStore.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.LastName, updatedUser.LastName)
	assert.Equal(t, user.Username, updatedUser.Username)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Password, updatedUser.Password)

	t.Run("Update user fail", func(t *testing.T) {
		user.ID = uuid.Nil
		err = testStorage.UserStore.UpdateUser(ctx, user)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, updatedUser.ID)
	})
}

func TestUserStore_DeleteUser(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	err = testStorage.UserStore.DeleteUser(ctx, user.ID)
	assert.NoError(t, err)

	deletedUser, err := testStorage.UserStore.GetUserByID(ctx, user.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, deletedUser)

	t.Run("Delete user fail", func(t *testing.T) {
		err = testStorage.UserStore.DeleteUser(ctx, user.ID)
		assert.Error(t, err)
	})
}

func TestUserStore_GetUserByUserID(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existUser)
	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, existUser.ID)
	})
}

func TestUserStore_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername(ctx, "test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, existUser.ID)
	})
}

func TestUserStore_GetUsersByUsername(t *testing.T) {
	ctx := context.Background()
	userJohn := createTestUser(t, "test", "test", "john", "john@test.com", "test")
	userJohnny := createTestUser(t, "test", "test", "johnny", "johnny@test.com", "test")
	userTest := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, userJohn)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(ctx, userJohnny)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(ctx, userTest)
	assert.NoError(t, err)

	t.Run("Get users by username success", func(t *testing.T) {
		users, err := testStorage.UserStore.GetUsersByUsername(ctx, "john")
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Equal(t, 2, len(users))
		assert.Contains(t, users[0].Username, "john")
		assert.Contains(t, users[1].Username, "john")

	})
	t.Run("Get users by username fail", func(t *testing.T) {
		users, err := testStorage.UserStore.GetUsersByUsername(ctx, "z")
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, users)
	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, userJohn.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, userJohnny.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, userTest.ID)

	})
}

func TestUserStore_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByEmail(ctx, "test@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, existUser.ID)
	})
}

func TestFollowStore_FollowUser(t *testing.T) {
	ctx := context.Background()
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(ctx, user1)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(ctx, user2)
	assert.NoError(t, err)

	follow := createTestFollow(t, user1.ID, user2.ID)

	err = testStorage.FollowStore.FollowUser(ctx, *follow)
	assert.NoError(t, err)

	t.Run("Follow user fail", func(t *testing.T) {
		err = testStorage.FollowStore.FollowUser(ctx, *follow)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.FollowStore.UnFollowUser(ctx, *follow)
		_ = testStorage.UserStore.DeleteUser(ctx, user1.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user2.ID)

	})

}

func TestFollowStore_UnFollowStore(t *testing.T) {
	ctx := context.Background()
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(ctx, user1)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(ctx, user2)
	assert.NoError(t, err)

	follow := createTestFollow(t, user1.ID, user2.ID)
	err = testStorage.FollowStore.FollowUser(ctx, *follow)
	assert.NoError(t, err)

	t.Run("Unfollow user success", func(t *testing.T) {
		err = testStorage.FollowStore.UnFollowUser(ctx, *follow)
		assert.NoError(t, err)
	})

	t.Run("Unfollow user fail", func(t *testing.T) {
		err = testStorage.FollowStore.UnFollowUser(ctx, *follow)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user1.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user2.ID)

	})

}

func TestFollowStore_GetFollowerByUserID(t *testing.T) {
	ctx := context.Background()
	user1 := createTestUser(t, "test1", "test1", "test1", "test1@test.com", "test1")
	user2 := createTestUser(t, "test2", "test2", "test2", "test2@test.com", "test2")

	err := testStorage.UserStore.CreateUser(ctx, user1)
	assert.NoError(t, err)
	err = testStorage.UserStore.CreateUser(ctx, user2)
	assert.NoError(t, err)

	follow := createTestFollow(t, user1.ID, user2.ID)
	err = testStorage.FollowStore.FollowUser(ctx, *follow)
	assert.NoError(t, err)

	followings, err := testStorage.FollowStore.GetFollowerByUserID(ctx, user2.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, followings)
	assert.Equal(t, 1, len(followings))
	assert.Equal(t, user1.ID, followings[0].UserID)
	assert.Equal(t, user2.ID, followings[0].FollowID)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user1.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user2.ID)
		_ = testStorage.FollowStore.UnFollowUser(ctx, *follow)

	})
}

func TestFollowStore_GetFollowinsByUserID(t *testing.T) {
	ctx := context.Background()
	user1 := createTestUser(t, "test1", "test1", "test1", "test1@test.com", "test1")
	user2 := createTestUser(t, "test2", "test2", "test2", "test2@test.com", "test2")

	err := testStorage.UserStore.CreateUser(ctx, user1)
	assert.NoError(t, err)
	err = testStorage.UserStore.CreateUser(ctx, user2)
	assert.NoError(t, err)

	follow := createTestFollow(t, user1.ID, user2.ID)
	err = testStorage.FollowStore.FollowUser(ctx, *follow)
	assert.NoError(t, err)

	followings, err := testStorage.FollowStore.GetFollowingByUserID(ctx, user1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, followings)
	assert.Equal(t, 1, len(followings))
	assert.Equal(t, user1.ID, followings[0].UserID)
	assert.Equal(t, user2.ID, followings[0].FollowID)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user1.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user2.ID)
		_ = testStorage.FollowStore.UnFollowUser(ctx, *follow)

	})
}

func TestPostStore_CreatePost(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	existPost, err := testStorage.PostStore.GetPostByID(ctx, post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existPost)

	assert.Equal(t, post.ID.String(), existPost.ID.String())
	assert.Equal(t, post.UserID.String(), existPost.UserID.String())
	assert.Equal(t, post.Content, existPost.Content)

	t.Run("Create post fail", func(t *testing.T) {
		err = testStorage.PostStore.CreatePost(ctx, post)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})

}

func TestPostStore_UpdatePost(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)
	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	post.Content = "updated"

	err = testStorage.PostStore.UpdatePost(ctx, post)
	assert.NoError(t, err)

	updatedPost, err := testStorage.PostStore.GetPostByID(ctx, post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, post.Content, updatedPost.Content)

	t.Run("Update post fail", func(t *testing.T) {
		post.ID = uuid.Nil
		err = testStorage.PostStore.UpdatePost(ctx, post)
		assert.Error(t, err)
	})
	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, updatedPost.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})

}

func TestPostStore_DeletePost(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	err = testStorage.PostStore.DeletePost(ctx, post.ID)
	assert.NoError(t, err)

	deletedPost, err := testStorage.PostStore.GetPostByID(ctx, post.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, deletedPost)

	t.Run("Delete post fail", func(t *testing.T) {
		err = testStorage.PostStore.DeletePost(ctx, post.ID)
		assert.Error(t, err)

	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestPostStore_GetPostsByUserIDByLimit(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	pagination.Limit = 5
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), user.ID, "public")
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(ctx, post)
		assert.NoError(t, err)
	}

	fivePosts, err := testStorage.PostStore.GetPostsByUserID(ctx, user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(fivePosts))

	pagination.Limit = 10

	allPosts, err := testStorage.PostStore.GetPostsByUserID(ctx, user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})

}

func TestPostStore_GetPostsByUserIDByQuery(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "1")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), user.ID, "public")
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(ctx, post)
		assert.NoError(t, err)
	}

	allPosts, err := testStorage.PostStore.GetPostsByUserID(ctx, user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allPosts))

	search.Query = "post"
	allPosts, err = testStorage.PostStore.GetPostsByUserID(ctx, user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestPostStore_GetPostDetailsByID(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)
	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	existPost, err := testStorage.PostStore.GetPostDetailsByID(ctx, post.ID, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existPost)
	assert.Empty(t, existPost.Comments)
	assert.Equal(t, 0, existPost.LikeCount)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestPostStore_HasAccessToPost(t *testing.T) {
	ctx := context.Background()
	user1 := createTestUser(t, "test1", "test1", "test1", "test1@test.com", "test1")
	user2 := createTestUser(t, "test2", "test2", "test2", "test2@test.com", "test2")

	err := testStorage.UserStore.CreateUser(ctx, user1)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(ctx, user2)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user1.ID, "private")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	hasAccess, err := testStorage.PostStore.HasAccessToPost(ctx, user2.ID, post.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, hasAccess)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user1.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user2.ID)
	})

}

func TestCommentStore_CreateComment(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(ctx, comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	t.Run("Create comment fail", func(t *testing.T) {
		err = testStorage.CommentStore.CreateComment(ctx, comment)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestCommentStore_UpdateComment(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(ctx, comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	comment.Content = "updated"

	err = testStorage.CommentStore.UpdateComment(ctx, comment)
	assert.NoError(t, err)

	updatedComment, err := testStorage.CommentStore.GetCommentByID(ctx, comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedComment)

	assert.Equal(t, comment.Content, updatedComment.Content)

	t.Run("Update comment fail", func(t *testing.T) {
		comment.ID = uuid.Nil
		err = testStorage.CommentStore.UpdateComment(ctx, comment)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.CommentStore.DeleteComment(ctx, updatedComment.ID)

		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestCommentStore_DeleteComment(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(ctx, comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	err = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
	assert.NoError(t, err)

	deletedComment, err := testStorage.CommentStore.GetCommentByID(ctx, comment.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, deletedComment)

	t.Run("Delete comment fail", func(t *testing.T) {
		err = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestLikeStore_LikePost(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	like := createTestLikePost(t, post.ID, user.ID)

	err = testStorage.LikeStore.LikePost(ctx, like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(ctx, post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Run("Like post fail", func(t *testing.T) {
		err = testStorage.PostStore.CreatePost(ctx, post)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.LikeStore.UnlikePost(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestLikeStore_UnlikePost(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	like := createTestLikePost(t, post.ID, user.ID)

	err = testStorage.LikeStore.LikePost(ctx, like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(ctx, post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikePost(ctx, post.ID, user.ID)
	assert.NoError(t, err)

	isLiked, err = testStorage.LikeStore.IsPostLiked(ctx, post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Run("Unlike post fail", func(t *testing.T) {
		err = testStorage.LikeStore.UnlikePost(ctx, post.ID, user.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestLikeStore_LikeComment(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	like := createTestLikeComment(t, comment.ID, user.ID)

	err = testStorage.LikeStore.LikeComment(ctx, like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(ctx, comment.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Run("Like post fail", func(t *testing.T) {
		err = testStorage.LikeStore.LikeComment(ctx, like)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.LikeStore.UnlikeComment(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestLikeStore_UnlikeComment(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	like := createTestLikeComment(t, comment.ID, user.ID)

	err = testStorage.LikeStore.LikeComment(ctx, like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(ctx, comment.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikeComment(ctx, comment.ID, user.ID)
	assert.NoError(t, err)
	isLiked, err = testStorage.LikeStore.IsCommentLiked(ctx, post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Run("Unlike comment fail", func(t *testing.T) {
		err = testStorage.LikeStore.UnlikeComment(ctx, comment.ID, user.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_CreateReply(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	assert.Equal(t, reply.ID.String(), existReply.ID.String())
	assert.Equal(t, reply.CommentID.String(), existReply.CommentID.String())
	assert.Equal(t, reply.UserID.String(), existReply.UserID.String())

	t.Run("Create reply fail", func(t *testing.T) {
		err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.LikeStore.UnlikeComment(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_CreateNestedReply(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	nestedReply := createTestNestedReply(t, reply.ID, user.ID, "nested_reply_test")
	err = testStorage.ReplyStore.CreateNestedReply(ctx, nestedReply)
	assert.NoError(t, err)

	existNestedReply, err := testStorage.ReplyStore.GetReplyByID(ctx, nestedReply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existNestedReply)

	assert.Equal(t, nestedReply.UserID.String(), existNestedReply.UserID.String())
	assert.Equal(t, nestedReply.Message, existNestedReply.Message)

	t.Run("Create nested reply fail", func(t *testing.T) {
		err = testStorage.ReplyStore.CreateNestedReply(ctx, nestedReply)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, reply.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, nestedReply.ID)
		_ = testStorage.LikeStore.UnlikeComment(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_UpdateReply(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	reply.Message = "updated"

	err = testStorage.ReplyStore.UpdateReply(ctx, reply)
	assert.NoError(t, err)

	updatedReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.Equal(t, reply.Message, updatedReply.Message)

	t.Run("Update reply fail", func(t *testing.T) {
		reply.ID = uuid.Nil
		err = testStorage.ReplyStore.UpdateReply(ctx, reply)
		assert.NoError(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.LikeStore.UnlikeComment(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_DeleteReply(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	err = testStorage.ReplyStore.DeleteReply(ctx, existReply.ID)
	assert.NoError(t, err)

	deletedReply, err := testStorage.ReplyStore.GetReplyByID(ctx, existReply.ID)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, deletedReply)

	t.Run("Delete reply fail", func(t *testing.T) {
		err = testStorage.ReplyStore.DeleteReply(ctx, existReply.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.LikeStore.UnlikeComment(ctx, post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_LikeReply(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	replyLike := createTestLikeReply(t, reply.ID, user.ID)

	err = testStorage.LikeStore.LikeReply(ctx, replyLike)
	assert.NoError(t, err)

	isLiked, err := testStorage.LikeStore.IsReplyLiked(ctx, reply.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Run("Like Reply fail", func(t *testing.T) {
		err = testStorage.LikeStore.LikeReply(ctx, replyLike)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_UnlikeReply(t *testing.T) {
	ctx := context.Background()

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	replyLike := createTestLikeReply(t, reply.ID, user.ID)

	err = testStorage.LikeStore.LikeReply(ctx, replyLike)
	assert.NoError(t, err)

	err = testStorage.LikeStore.UnlikeReply(ctx, reply.ID, user.ID)
	assert.NoError(t, err)

	t.Run("Unlike reply fail", func(t *testing.T) {
		err = testStorage.LikeStore.UnlikeReply(ctx, reply.ID, user.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_GetRepliesByCommentID(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(ctx, reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	commentReplies, err := testStorage.ReplyStore.GetRepliesByCommentID(ctx, comment.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(commentReplies))
	assert.Equal(t, existReply.ID, commentReplies[0].ID)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestReplyStore_GetRepliesByParentID(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID, "public")

	err = testStorage.PostStore.CreatePost(ctx, post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(ctx, comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(ctx, reply)
	assert.NoError(t, err)

	nestedReply := createTestNestedReply(t, reply.ID, user.ID, "nested")
	err = testStorage.ReplyStore.CreateNestedReply(ctx, nestedReply)

	assert.NoError(t, err)

	replies, err := testStorage.ReplyStore.GetRepliesByParentID(ctx, nestedReply.ParentID, user.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, replies)
	assert.Equal(t, 1, len(replies))
	assert.Equal(t, nestedReply.ID, replies[0].ID)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(ctx, post.ID)
		_ = testStorage.CommentStore.DeleteComment(ctx, comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(ctx, post.ID)
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
	})
}

func TestTagStore_AddPostTags(t *testing.T) {
	ctx := context.Background()
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(ctx, user)
	assert.NoError(t, err)

	post1 := createTestPost(t, "test", user.ID, "public")
	err = testStorage.PostStore.CreatePost(ctx, post1)
	assert.NoError(t, err)

	post2 := createTestPost(t, "test", user.ID, "public")
	err = testStorage.PostStore.CreatePost(ctx, post2)
	assert.NoError(t, err)

	tags := createTestTags("test1", "test2")
	tagStrs := []string{"test1", "test2"}

	err = testStorage.TagStore.CreateMultipleTags(ctx, tags)
	assert.NoError(t, err)

	err = testStorage.TagStore.CreatePostTag(ctx, tagStrs, post1.ID)
	assert.NoError(t, err)

	err = testStorage.TagStore.CreatePostTag(ctx, []string{"test1"}, post2.ID)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	posts, err := testStorage.PostStore.GetPostsByTag(ctx, pagination, "test1", user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	posts, err = testStorage.PostStore.GetPostsByTag(ctx, pagination, "test2", user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))

	t.Run("Create post tag fail", func(t *testing.T) {
		err = testStorage.TagStore.CreatePostTag(ctx, nil, post2.ID)
		assert.Error(t, err)
	})

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(ctx, user.ID)
		_ = testStorage.PostStore.DeletePost(ctx, post1.ID)
		_ = testStorage.PostStore.DeletePost(ctx, post2.ID)
	})
}

func TestMain(m *testing.M) {
	testStorage = NewPostgresTestStorage()
	testDB = testStorage.UserStore.(*UserStore).DB
	cleanupAllTables()
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}
