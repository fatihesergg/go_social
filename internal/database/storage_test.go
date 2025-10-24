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
		ReplyStore:   NewReplyStore(db),
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

func createTestPost(t *testing.T, content string, userID uuid.UUID) *model.Post {
	t.Helper()
	return &model.Post{
		ID:      uuid.New(),
		Content: content,
		UserID:  userID,
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

func TestUserStore_CreateUser(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByID(user.ID)
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

	user.Name = "test_update"
	user.LastName = "test_update"
	user.Username = "test_update"
	user.Email = "test_update@test.com"
	user.Password = "test_update"

	err = testStorage.UserStore.UpdateUser(user)
	assert.NoError(t, err)

	updatedUser, err := testStorage.UserStore.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, user.Name, updatedUser.Name)
	assert.Equal(t, user.LastName, updatedUser.LastName)
	assert.Equal(t, user.Username, updatedUser.Username)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Password, updatedUser.Password)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(updatedUser.ID)
	})
}

func TestUserStore_DeleteUser(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	err = testStorage.UserStore.DeleteUser(user.ID)
	assert.NoError(t, err)

	deletedUser, err := testStorage.UserStore.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserStore_GetUserByUserID(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	existUser, err := testStorage.UserStore.GetUserByID(user.ID)
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

	follow := createTestFollow(t, user1.ID, user2.ID)

	err = testStorage.FollowStore.FollowUser(*follow)
	assert.NoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(follows))
	first := follows[0]
	assert.Equal(t, first.UserID.String(), user1.ID.String())
	assert.Equal(t, first.FollowID.String(), user2.ID.String())

	t.Cleanup(func() {
		_ = testStorage.FollowStore.UnFollowUser(*follow)
		_ = testStorage.UserStore.DeleteUser(user1.ID)
		_ = testStorage.UserStore.DeleteUser(user2.ID)

	})

}

func TestFollowStore_UnFollowStore(t *testing.T) {
	user1 := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	user2 := createTestUser(t, "test_2", "test_2", "test_2", "test_2@test.com", "test_2")

	err := testStorage.UserStore.CreateUser(user1)
	assert.NoError(t, err)

	err = testStorage.UserStore.CreateUser(user2)
	assert.NoError(t, err)

	follow := createTestFollow(t, user1.ID, user2.ID)
	err = testStorage.FollowStore.FollowUser(*follow)
	assert.NoError(t, err)

	err = testStorage.FollowStore.UnFollowUser(*follow)
	assert.NoError(t, err)

	follows, err := testStorage.FollowStore.GetFollowingByUserID(user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(follows))

	t.Cleanup(func() {
		_ = testStorage.FollowStore.UnFollowUser(*follow)
		_ = testStorage.UserStore.DeleteUser(user1.ID)
		_ = testStorage.UserStore.DeleteUser(user2.ID)

	})

}

func TestPostStore_CreatePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	existPost, err := testStorage.PostStore.GetPostByID(post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existPost)

	assert.Equal(t, post.ID.String(), existPost.ID.String())
	assert.Equal(t, post.UserID.String(), existPost.UserID.String())
	assert.Equal(t, post.Content, existPost.Content)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})

}

func TestPostStore_UpdatePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)
	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	post.Content = "updated"

	err = testStorage.PostStore.UpdatePost(post)
	assert.NoError(t, err)

	updatedPost, err := testStorage.PostStore.GetPostByID(post.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, post.Content, updatedPost.Content)

	t.Cleanup(func() {
		_ = testStorage.PostStore.DeletePost(updatedPost.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})

}

func TestPostStore_DeletePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	err = testStorage.PostStore.DeletePost(post.ID)
	assert.NoError(t, err)

	deletedPost, err := testStorage.PostStore.GetPostByID(post.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedPost)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestPostStore_GetPostsByUserIDByLimit(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	pagination.Limit = 5
	search := createTestSearch(t, "")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), user.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		assert.NoError(t, err)
	}

	fivePosts, err := testStorage.PostStore.GetPostsByUserID(user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(fivePosts))

	pagination.Limit = 10

	allPosts, err := testStorage.PostStore.GetPostsByUserID(user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})

}

func TestPostStore_GetPostsByUserIDByQuery(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")
	pagination := createTestPagination(t)
	search := createTestSearch(t, "1")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	var posts []*model.Post

	for i := 1; i < 11; i++ {
		post := createTestPost(t, fmt.Sprintf("test_post_%d", i), user.ID)
		posts = append(posts, post)

	}

	for _, post := range posts {

		err = testStorage.PostStore.CreatePost(post)
		assert.NoError(t, err)
	}

	allPosts, err := testStorage.PostStore.GetPostsByUserID(user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(allPosts))

	search.Query = "post"
	allPosts, err = testStorage.PostStore.GetPostsByUserID(user.ID, pagination, search)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(allPosts))

	t.Cleanup(func() {
		for _, post := range posts {
			_ = testStorage.PostStore.DeletePost(post.ID)
		}
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestCommentStore_CreateComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	t.Cleanup(func() {
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestCommentStore_UpdateComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	comment.Content = "updated"

	err = testStorage.CommentStore.UpdateComment(comment)
	assert.NoError(t, err)

	updatedComment, err := testStorage.CommentStore.GetCommentByID(comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedComment)

	assert.Equal(t, comment.Content, updatedComment.Content)

	t.Cleanup(func() {

		_ = testStorage.CommentStore.DeleteComment(updatedComment.ID)

		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestCommentStore_DeleteComment(t *testing.T) {

	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)
	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	existComment, err := testStorage.CommentStore.GetCommentByID(comment.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existComment)

	assert.Equal(t, comment.Content, existComment.Content)
	assert.Equal(t, comment.PostID.String(), existComment.PostID.String())
	assert.Equal(t, user.ID.String(), existComment.UserID.String())

	err = testStorage.CommentStore.DeleteComment(comment.ID)
	assert.NoError(t, err)

	deletedComment, err := testStorage.CommentStore.GetCommentByID(comment.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedComment)

	t.Cleanup(func() {
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestLikeStore_LikePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	like := createTestLikePost(t, post.ID, user.ID)

	err = testStorage.LikeStore.LikePost(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.LikeStore.UnlikePost(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestLikeStore_UnlikePost(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	like := createTestLikePost(t, post.ID, user.ID)

	err = testStorage.LikeStore.LikePost(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsPostLiked(post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikePost(post.ID, user.ID)
	assert.NoError(t, err)

	isLiked, err = testStorage.LikeStore.IsPostLiked(post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestLikeStore_LikeComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	like := createTestLikeComment(t, comment.ID, user.ID)

	err = testStorage.LikeStore.LikeComment(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(comment.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.LikeStore.UnlikeComment(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestLikeStore_UnlikeComment(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	like := createTestLikeComment(t, comment.ID, user.ID)

	err = testStorage.LikeStore.LikeComment(like)
	assert.NoError(t, err)
	isLiked, err := testStorage.LikeStore.IsCommentLiked(comment.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, true, isLiked)

	err = testStorage.LikeStore.UnlikeComment(comment.ID, user.ID)
	assert.NoError(t, err)
	isLiked, err = testStorage.LikeStore.IsCommentLiked(post.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, false, isLiked)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestReplyStore_CreateReply(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	assert.Equal(t, reply.ID.String(), existReply.ID.String())
	assert.Equal(t, reply.CommentID.String(), existReply.CommentID.String())
	assert.Equal(t, reply.UserID.String(), existReply.UserID.String())

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(post.ID)
		_ = testStorage.LikeStore.UnlikeComment(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestReplyStore_CreateNestedReply(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	nestedReply := createTestNestedReply(t, reply.ID, user.ID, "nested_reply_test")
	err = testStorage.ReplyStore.CreateNestedReply(nestedReply)
	assert.NoError(t, err)

	existNestedReply, err := testStorage.ReplyStore.GetReplyByID(nestedReply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existNestedReply)

	assert.Equal(t, nestedReply.UserID.String(), existNestedReply.UserID.String())
	assert.Equal(t, nestedReply.Message, existNestedReply.Message)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(reply.ID)
		_ = testStorage.ReplyStore.DeleteReply(nestedReply.ID)
		_ = testStorage.LikeStore.UnlikeComment(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestReplyStore_UpdateReply(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(reply)

	assert.NoError(t, err)

	reply.Message = "updated"

	err = testStorage.ReplyStore.UpdateReply(reply)
	assert.NoError(t, err)

	updatedReply, err := testStorage.ReplyStore.GetReplyByID(reply.ID)
	assert.NoError(t, err)
	assert.Equal(t, reply.Message, updatedReply.Message)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(post.ID)
		_ = testStorage.LikeStore.UnlikeComment(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
	})
}

func TestReplyStore_DeleteReply(t *testing.T) {
	user := createTestUser(t, "test", "test", "test", "test@test.com", "test")

	err := testStorage.UserStore.CreateUser(user)
	assert.NoError(t, err)

	post := createTestPost(t, "test", user.ID)

	err = testStorage.PostStore.CreatePost(post)
	assert.NoError(t, err)

	comment := createTestComment(t, "test", post.ID, user.ID)

	err = testStorage.CommentStore.CreateComment(comment)
	assert.NoError(t, err)

	reply := createTestCommentReply(t, comment.ID, user.ID, "test")

	err = testStorage.ReplyStore.CreateCommentReply(reply)

	assert.NoError(t, err)

	existReply, err := testStorage.ReplyStore.GetReplyByID(reply.ID)
	assert.NoError(t, err)
	assert.NotNil(t, existReply)

	err = testStorage.ReplyStore.DeleteReply(existReply.ID)
	assert.NoError(t, err)

	deletedReply, err := testStorage.ReplyStore.GetReplyByID(existReply.ID)
	assert.NoError(t, err)
	assert.Nil(t, deletedReply)

	t.Cleanup(func() {

		_ = testStorage.PostStore.DeletePost(post.ID)
		_ = testStorage.CommentStore.DeleteComment(comment.ID)
		_ = testStorage.ReplyStore.DeleteReply(post.ID)
		_ = testStorage.LikeStore.UnlikeComment(post.ID, user.ID)
		_ = testStorage.UserStore.DeleteUser(user.ID)
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
