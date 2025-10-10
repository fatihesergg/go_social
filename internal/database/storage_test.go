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
		Name:     name,
		LastName: lastName,
		Username: username,
		Email:    email,
		Password: password,
	}
}

func createTestPost(t *testing.T, content string, userID uuid.UUID) *model.Post {
	t.Helper()
	return &model.Post{
		Content: content,
		UserID:  userID,
	}
}

func createTestComment(t *testing.T, content string, postID, userID uuid.UUID) *model.Comment {
	t.Helper()

	return &model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}
}

func createTestLikePost(t *testing.T, postID, userID uuid.UUID) *model.PostLike {
	t.Helper()
	return &model.PostLike{
		PostID: postID,
		UserID: userID,
	}
}
func createTestLikeComment(t *testing.T, commentID, userID uuid.UUID) *model.CommentLike {
	t.Helper()
	return &model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
}
func TestUserStore_CreateUser(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)
	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_UpdateUser(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	existUser.Name = "test_update"
	existUser.LastName = "test_update"
	existUser.Username = "test_update"
	existUser.Email = "test_update@test.com"
	existUser.Password = "test_update"

	err = testStorage.UserStore.UpdateUser(existUser)
	assert.NoError(t, err)

	updatedUser, err := testStorage.UserStore.GetUserByUsername("test_update")
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, existUser.Name, updatedUser.Name)
	assert.Equal(t, existUser.LastName, updatedUser.LastName)
	assert.Equal(t, existUser.Username, updatedUser.Username)
	assert.Equal(t, existUser.Email, updatedUser.Email)
	assert.Equal(t, existUser.Password, updatedUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(updatedUser.ID)
	})
}

func TestUserStore_DeleteUser(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)

	err = testStorage.UserStore.DeleteUser(existUser.ID)
	assert.NoError(t, err)

	deletedUser, err := testStorage.UserStore.GetUserByUsername(existUser.Username)
	assert.NoError(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserStore_GetUserByUserID(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	existUser, err = testStorage.UserStore.GetUserByID(existUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existUser)
	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_GetUserByUsername(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestUserStore_GetUserByEmail(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByEmail("test@test.com")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	assert.Equal(t, user.Name, existUser.Name)
	assert.Equal(t, user.LastName, existUser.LastName)
	assert.Equal(t, user.Username, existUser.Username)
	assert.Equal(t, user.Email, existUser.Email)
	assert.Equal(t, user.Password, existUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestFollowStore_FollowUser(t *testing.T) {
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(user1)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(user2)
	assert.NoError(t, err)

	existUser1, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser1)
	existUser2, err := testStorage.UserStore.GetUserByUsername("test_2")
	assert.NoError(t, err)
	assert.NotNil(t, existUser2)

	err = testStorage.FollowStore.FollowUser(existUser1.ID, existUser2.ID)
	assert.NoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(follows))
	first := follows[0]
	assert.Equal(t, first.UserID.String(), existUser1.ID.String())
	assert.Equal(t, first.FollowID.String(), existUser2.ID.String())

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
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(user2)
	assert.NoError(t, err)

	existUser1, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser1)
	existUser2, err := testStorage.UserStore.GetUserByUsername("test_2")
	assert.NoError(t, err)
	assert.NotNil(t, existUser2)

	err = testStorage.FollowStore.FollowUser(existUser1.ID, existUser2.ID)
	assert.NoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(follows))

	err = testStorage.FollowStore.UnFollowUser(existUser1.ID, existUser2.ID)
	assert.NoError(t, err)

	follows, err = testStorage.FollowStore.GetFollowingByUserID(existUser1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(follows))

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))
	first := existPosts[0]

	assert.Equal(t, existUser.ID.String(), first.User.ID.String())
	assert.Equal(t, post.Content, first.Content)

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))
	first := existPosts[0]

	assert.Equal(t, existUser.ID.String(), first.User.ID.String())
	assert.Equal(t, post.Content, first.Content)

	first.Content = "updated"

	err = testStorage.PostStore.UpdatePost(&first)
	assert.NoError(t, err)

	updatedPost, err := testStorage.PostStore.GetPostDetailsByID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, first.Content, updatedPost.Content)

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	existPosts, err := testStorage.PostStore.GetPostsByUserID(post.UserID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))
	first := existPosts[0]

	assert.Equal(t, existUser.ID.String(), first.User.ID.String())
	assert.Equal(t, post.Content, first.Content)

	err = testStorage.PostStore.DeletePost(first.ID)
	assert.NoError(t, err)

	deletedPost, err := testStorage.PostStore.GetPostDetailsByID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedPost)

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), existUser.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		assert.NoError(t, err)
	}

	fivePosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(fivePosts))

	pagination.Limit = 10

	allPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), existUser.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		assert.NoError(t, err)
	}

	allPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allPosts))

	search.Query = "post"
	allPosts, err = testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(comments))

	firstComment := comments[0]
	assert.Equal(t, comment.Content, firstComment.Content)
	assert.Equal(t, comment.PostID.String(), firstComment.PostID.String())
	assert.Equal(t, existUser.ID.String(), firstComment.User.ID.String())

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
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(comments))

	firstComment := comments[0]
	assert.Equal(t, comment.Content, firstComment.Content)
	assert.Equal(t, comment.PostID.String(), firstComment.PostID.String())
	assert.Equal(t, existUser.ID.String(), firstComment.User.ID.String())

	firstComment.Content = "updated"

	err = testStorage.CommentStore.UpdateComment(&firstComment)
	assert.NoError(t, err)

	updatedComment, err := testStorage.CommentStore.GetCommentByID(firstComment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedComment)

	assert.Equal(t, firstComment.Content, updatedComment.Content)

	t.Cleanup(func() {

		_ = testStorage.CommentStore.DeleteComment(updatedComment.ID)

		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestCommentStore_DeleteComment(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	comments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(comments))

	firstComment := comments[0]
	assert.Equal(t, comment.Content, firstComment.Content)
	assert.Equal(t, comment.PostID.String(), firstComment.PostID.String())
	assert.Equal(t, existUser.ID.String(), firstComment.User.ID.String())

	err = testStorage.CommentStore.DeleteComment(firstComment.ID)
	assert.NoError(t, err)

	deletedComment, err := testStorage.CommentStore.GetCommentByID(firstComment.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedComment)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestLikeStore_LikePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	like := createTestLikePost(t, first.ID, existUser.ID)

	err = testStorage.LikeStore.LikePost(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(first.ID)
		_ = testStorage.LikeStore.UnlikePost(first.ID, existUser.ID)
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestLikeStore_UnlikePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	like := createTestLikePost(t, first.ID, existUser.ID)

	err = testStorage.LikeStore.LikePost(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikePost(first.ID, existUser.ID)
	assert.NoError(t, err)
	isLiked, err = testStorage.LikeStore.IsPostLiked(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(first.ID)
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestLikeStore_LikeComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	existComments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existComments))
	firstComment := existComments[0]

	like := createTestLikeComment(t, firstComment.ID, existUser.ID)

	err = testStorage.LikeStore.LikeComment(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(firstComment.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(first.ID)
		_ = testStorage.CommentStore.DeleteComment(firstComment.ID)
		_ = testStorage.LikeStore.UnlikeComment(first.ID, existUser.ID)
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
	})
}

func TestLikeStore_UnlikeComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByUsername("test")
	assert.NoError(t, err)
	assert.NotNil(t, existUser)

	post := createTestPost(t, "test", existUser.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	pagination := createTestPagination(t)
	search := createTestSearch(t, "")
	existPosts, err := testStorage.PostStore.GetPostsByUserID(existUser.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existPosts))

	first := existPosts[0]

	comment := createTestComment(t, "test", first.ID, existUser.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	existComments, err := testStorage.CommentStore.GetCommentsByPostID(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(existComments))
	firstComment := existComments[0]

	like := createTestLikeComment(t, firstComment.ID, existUser.ID)

	err = testStorage.LikeStore.LikeComment(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(firstComment.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikeComment(firstComment.ID, existUser.ID)
	assert.NoError(t, err)
	isLiked, err = testStorage.LikeStore.IsCommentLiked(first.ID, existUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(first.ID)
		_ = testStorage.UserStore.DeleteUser(existUser.ID)
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
