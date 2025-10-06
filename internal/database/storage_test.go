package database

import (
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
)

var testDB *sql.DB
var testStorage *Storage

func NewInMemoryStorageForTest() *Storage {
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
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func AssertNil(t *testing.T, obj any) {
	// nil interface !!!

	switch v := obj.(type) {
	case *model.User:
		if v != nil {
			t.Errorf("Expected nil, but got: %v", v)
		}
	}
}

func AssertNotNil(t *testing.T, obj any) {
	if obj == nil {
		t.Errorf("Expected not nil, but got nil")
	}
}

func AssertStringEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected string %s, got %s", expected, actual)
	}
}

func AssertIntEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected int %d, got %d", expected, actual)
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
	return Pagination{
		Limit:  20,
		Offset: 0,
	}
}

func createTestSearch(t *testing.T, query string) Search {
	return Search{
		Query: query,
	}
}

func createTestUser(t *testing.T, name, lastName, username, email, password string) *model.User {
	return &model.User{
		Name:     name,
		LastName: lastName,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func createTestPost(t *testing.T, content string, userID uuid.UUID) *model.Post {
	return &model.Post{
		Content: content,
		UserID:  userID,
	}
}

func createTestComment(t *testing.T, content string, postID, userID uuid.UUID) *model.Comment {

	return &model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
}

func TestUserStore_CreateUser(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)
	AssertStringEqual(t, user.Name, existUser.Name)
	AssertStringEqual(t, user.LastName, existUser.LastName)
	AssertStringEqual(t, user.Username, existUser.Username)
	AssertStringEqual(t, user.Email, existUser.Email)
	AssertStringEqual(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_UpdateUser(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	existUser.Name = "test_update"
	existUser.LastName = "test_update"
	existUser.Username = "test_update"
	existUser.Email = "test_update@test.com"
	existUser.Password = "test_update"

	err = testStorage.UserStore.UpdateUser(existUser)
	AssertNoError(t, err)

	updatedUser, err := testStorage.UserStore.GetUserByUsername("test_update")
	AssertNoError(t, err)
	AssertNotNil(t, updatedUser)
	AssertStringEqual(t, existUser.Name, updatedUser.Name)
	AssertStringEqual(t, existUser.LastName, updatedUser.LastName)
	AssertStringEqual(t, existUser.Username, updatedUser.Username)
	AssertStringEqual(t, existUser.Email, updatedUser.Email)
	AssertStringEqual(t, existUser.Password, updatedUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(updatedUser.ID)
	})
}

func TestUserStore_DeleteUser(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)

	err = testStorage.UserStore.DeleteUser(existUser.ID)
	AssertNoError(t, err)

	deletedUser, err := testStorage.UserStore.GetUserByUsername(existUser.Username)
	AssertNoError(t, err)
	AssertNil(t, deletedUser)
}

func TestUserStore_GetUserByUserID(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	existUser, err = testStorage.UserStore.GetUserByID(existUser.ID)
	AssertNoError(t, err)
	AssertNotNil(t, existUser)
	AssertStringEqual(t, user.Name, existUser.Name)
	AssertStringEqual(t, user.LastName, existUser.LastName)
	AssertStringEqual(t, user.Username, existUser.Username)
	AssertStringEqual(t, user.Email, existUser.Email)
	AssertStringEqual(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_GetUserByUsername(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	AssertStringEqual(t, user.Name, existUser.Name)
	AssertStringEqual(t, user.LastName, existUser.LastName)
	AssertStringEqual(t, user.Username, existUser.Username)
	AssertStringEqual(t, user.Email, existUser.Email)
	AssertStringEqual(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_GetUserByEmail(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByEmail("test@test.com")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	AssertStringEqual(t, user.Name, existUser.Name)
	AssertStringEqual(t, user.LastName, existUser.LastName)
	AssertStringEqual(t, user.Username, existUser.Username)
	AssertStringEqual(t, user.Email, existUser.Email)
	AssertStringEqual(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestFollowStore_FollowUser(t *testing.T) {
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(user1)
	AssertNoError(t, err)

	err = testStorage.UserStore.CreateUser(user2)
	AssertNoError(t, err)

	existUser1, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser1)
	existUser2, err := testStorage.UserStore.GetUserByUsername("test_2")
	AssertNoError(t, err)
	AssertNotNil(t, existUser2)

	err = testStorage.FollowStore.FollowUser(existUser1.ID, existUser2.ID)
	AssertNoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(follows))
	first := follows[0]
	AssertStringEqual(t, first.UserID.String(), existUser1.ID.String())
	AssertStringEqual(t, first.FollowID.String(), existUser2.ID.String())

	t.Cleanup(func() {
		_ = testStorage.FollowStore.UnFollowUser(existUser1.ID, existUser2.ID)
		_ = testStorage.UserStore.DeleteUser(existUser1.ID)
		_ = testStorage.UserStore.DeleteUser(existUser2.ID)

	})

}

func TestFollowStore_UnFollowStore(t *testing.T) {
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(user1)
	AssertNoError(t, err)

	err = testStorage.UserStore.CreateUser(user2)
	AssertNoError(t, err)

	existUser1, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser1)
	existUser2, err := testStorage.UserStore.GetUserByUsername("test_2")
	AssertNoError(t, err)
	AssertNotNil(t, existUser2)

	err = testStorage.FollowStore.FollowUser(existUser1.ID, existUser2.ID)
	AssertNoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(follows))

	err = testStorage.FollowStore.UnFollowUser(existUser1.ID, existUser2.ID)
	AssertNoError(t, err)

	follows, err = testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 0, len(follows))

	t.Cleanup(func() {
		_ = testStorage.FollowStore.UnFollowUser(existUser1.ID, existUser2.ID)
		_ = testStorage.UserStore.DeleteUser(existUser1.ID)
		_ = testStorage.UserStore.DeleteUser(existUser2.ID)

	})

}

func TestPostStore_CreatePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))
	first := existPosts[0]

	AssertStringEqual(t, existUser.ID.String(), first.User.ID.String())
	AssertStringEqual(t, post.Content, first.Content)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(first.ID)
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})

}

func TestPostStore_UpdatePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))
	first := existPosts[0]

	AssertStringEqual(t, existUser.ID.String(), first.User.ID.String())
	AssertStringEqual(t, post.Content, first.Content)

	first.Content = "updated"

	err = testStorage.PostStore.UpdatePost(&first)
	AssertNoError(t, err)

	updatedPost, err := testStorage.PostStore.GetPostByID(first.ID)
	AssertNoError(t, err)
	AssertNotNil(t, updatedPost)
	AssertStringEqual(t, first.Content, updatedPost.Content)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(updatedPost.ID)
		_ = testStorage.UserStore.DeleteUser(updatedPost.User.ID)
	})

}

func TestPostStore_DeletePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))
	first := existPosts[0]

	AssertStringEqual(t, existUser.ID.String(), first.User.ID.String())
	AssertStringEqual(t, post.Content, first.Content)

	err = testStorage.PostStore.DeletePost(first.ID)
	AssertNoError(t, err)

	deletedPost, err := testStorage.PostStore.GetPostByID(first.ID)
	AssertNoError(t, err)
	AssertNil(t, deletedPost)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestPostStore_GetPostsByUserIDByLimit(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	pagination.Limit = 5
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), existUser.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		AssertNoError(t, err)
	}

	fivePosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 5, len(fivePosts))

	pagination.Limit = 10

	allPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})

}

func TestPostStore_GetPostsByUserIDByQuery(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "1")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), existUser.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		AssertNoError(t, err)
	}

	allPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 2, len(allPosts))

	search.Query = "post"
	allPosts, err = testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestCommentStore_CreateComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	AssertNoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(comments))

	firstComment := comments[0]
	AssertStringEqual(t, comment.Content, firstComment.Content)
	AssertStringEqual(t, comment.PostID.String(), firstComment.PostID.String())
	AssertStringEqual(t, comment.UserID.String(), firstComment.UserID.String())

	t.Cleanup(func() {
		for _, comment := range comments {
			_ = testStorage.CommentStore.DeleteComment(comment.ID)
		}
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestCommentStore_UpdateComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	AssertNoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(comments))

	firstComment := comments[0]
	AssertStringEqual(t, comment.Content, firstComment.Content)
	AssertStringEqual(t, comment.PostID.String(), firstComment.PostID.String())
	AssertStringEqual(t, comment.UserID.String(), firstComment.UserID.String())

	firstComment.Content = "updated"

	err = testStorage.CommentStore.UpdateComment(&firstComment)
	AssertNoError(t, err)

	updatedComment, err := testStorage.CommentStore.GetCommentByID(firstComment.ID)
	AssertNoError(t, err)
	AssertNotNil(t, updatedComment)

	AssertStringEqual(t, firstComment.Content, updatedComment.Content)

	t.Cleanup(func() {

		_ = testStorage.CommentStore.DeleteComment(updatedComment.ID)

		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestCommentStore_DeleteComment(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	AssertNoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	AssertNoError(t, err)
	AssertNotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	AssertNoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	AssertNoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID)
	AssertNoError(t, err)
	AssertIntEqual(t, 1, len(comments))

	firstComment := comments[0]
	AssertStringEqual(t, comment.Content, firstComment.Content)
	AssertStringEqual(t, comment.PostID.String(), firstComment.PostID.String())
	AssertStringEqual(t, comment.UserID.String(), firstComment.UserID.String())

	err = testStorage.CommentStore.DeleteComment(firstComment.ID)
	AssertNoError(t, err)

	deletedComment, err := testStorage.CommentStore.GetCommentByID(firstComment.ID)
	AssertNoError(t, err)
	AssertNil(t, deletedComment)

}

func TestMain(m *testing.M) {
	testStorage = NewInMemoryStorageForTest()
	testDB = testStorage.UserStore.(*UserStore).DB
	cleanupAllTables()
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}
