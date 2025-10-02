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
	if obj != nil {
		t.Errorf("Expected nil, but got: %v", obj)
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

func cleanup() {
	tables := []string{"comment_likes", "post_likes", "follows", "comments", "posts", "users"}
	for _, table := range tables {
		fmt.Printf("Truncating table %s\n", table)
		_, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			fmt.Printf("Failed to truncate table %s: %v\n", table, err)
		}
	}
}

var testUser = model.User{
	Name:     "Test",
	LastName: "User",
	Username: "testuser",
	Email:    "testuser@example.com",
	Password: "password",
}

func TestUserRepo(t *testing.T) {

	t.Run("Create User", func(t *testing.T) {
		err := testStorage.UserStore.CreateUser(testUser)
		AssertNoError(t, err)
	})

	t.Run("Get User By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUser.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Email, user.Email)
		AssertStringEqual(t, testUser.Name, user.Name)
		AssertStringEqual(t, testUser.LastName, user.LastName)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Password, user.Password)

		// Set testUser ID for future tests
		testUser.ID = user.ID
	})

	t.Run("Delete User", func(t *testing.T) {
		err := testStorage.UserStore.DeleteUser(testUser.ID)
		AssertNoError(t, err)
	})

	t.Run("Get Deleted User By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUser.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
	})

	t.Run("Create User Again", func(t *testing.T) {
		err := testStorage.UserStore.CreateUser(testUser)
		AssertNoError(t, err)

	})

	t.Run("Get User By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUser.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Email, user.Email)
		AssertStringEqual(t, testUser.Name, user.Name)
		AssertStringEqual(t, testUser.LastName, user.LastName)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Password, user.Password)
	})

	t.Run("Get User By Email", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByEmail(testUser.Email)
		AssertNoError(t, err)
		AssertNotNil(t, user)

		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Email, user.Email)
		AssertStringEqual(t, testUser.Name, user.Name)
		AssertStringEqual(t, testUser.LastName, user.LastName)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Password, user.Password)
		testUser.ID = user.ID
	})

	t.Run("Get User By ID", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByID(testUser.ID)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUser.ID.String(), user.ID.String())
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Email, user.Email)
		AssertStringEqual(t, testUser.Name, user.Name)
		AssertStringEqual(t, testUser.LastName, user.LastName)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Password, user.Password)
	})

	t.Run("Update User", func(t *testing.T) {
		testUser.Name = "UpdatedName"
		testUser.LastName = "UpdatedLastName"
		testUser.Username = "updatedusername"
		testUser.Email = "updateduser@example.com"
		err := testStorage.UserStore.UpdateUser(testUser)
		AssertNoError(t, err)
	})

	t.Run("Get Updated User By ID", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByID(testUser.ID)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUser.ID.String(), user.ID.String())
		AssertStringEqual(t, testUser.Name, user.Name)
		AssertStringEqual(t, testUser.LastName, user.LastName)
		AssertStringEqual(t, testUser.Username, user.Username)
		AssertStringEqual(t, testUser.Email, user.Email)
		AssertStringEqual(t, testUser.Password, user.Password)
	})
}

func TestMain(m *testing.M) {
	testStorage = NewInMemoryStorageForTest()
	testDB = testStorage.UserStore.(*UserStore).DB
	cleanup()
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}
