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

func AssertIntEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected int %d, got %d", expected, actual)
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

var testUserJohn = model.User{
	Name:     "john",
	LastName: "doe",
	Username: "john_doe",
	Email:    "john_doe@example.com",
	Password: "password",
}

var testUserAlice = model.User{
	Name:     "alice",
	LastName: "doe",
	Username: "alice_doe",
	Email:    "alice_doe@example.com",
	Password: "password",
}

func TestUserRepo(t *testing.T) {

	t.Run("Create User testUserJohn", func(t *testing.T) {
		err := testStorage.UserStore.CreateUser(testUserJohn)
		AssertNoError(t, err)
	})

	t.Run("Create User testUserAlice", func(t *testing.T) {
		err := testStorage.UserStore.CreateUser(testUserAlice)
		AssertNoError(t, err)
	})

	t.Run("Get testUserJohn By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUserJohn.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Email, user.Email)
		AssertStringEqual(t, testUserJohn.Name, user.Name)
		AssertStringEqual(t, testUserJohn.LastName, user.LastName)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Password, user.Password)

		// Set testUserJohn ID for future tests
		testUserJohn.ID = user.ID
	})

	t.Run("Get testUserAlice By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUserAlice.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUserAlice.Username, user.Username)
		AssertStringEqual(t, testUserAlice.Email, user.Email)
		AssertStringEqual(t, testUserAlice.Name, user.Name)
		AssertStringEqual(t, testUserAlice.LastName, user.LastName)
		AssertStringEqual(t, testUserAlice.Username, user.Username)
		AssertStringEqual(t, testUserAlice.Password, user.Password)

		// Set testUserAlice ID for future tests
		testUserAlice.ID = user.ID
	})

	t.Run("Delete testUserJohn", func(t *testing.T) {
		err := testStorage.UserStore.DeleteUser(testUserJohn.ID)
		AssertNoError(t, err)
	})

	t.Run("Get Deleted testUserJohn By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUserJohn.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
	})

	t.Run("Create testUserJohn again", func(t *testing.T) {
		err := testStorage.UserStore.CreateUser(testUserJohn)
		AssertNoError(t, err)

	})

	t.Run("Get testUserJohn By Username", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByUsername(testUserJohn.Username)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Email, user.Email)
		AssertStringEqual(t, testUserJohn.Name, user.Name)
		AssertStringEqual(t, testUserJohn.LastName, user.LastName)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Password, user.Password)
	})

	t.Run("Get testUserJohn By Email", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByEmail(testUserJohn.Email)
		AssertNoError(t, err)
		AssertNotNil(t, user)

		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Email, user.Email)
		AssertStringEqual(t, testUserJohn.Name, user.Name)
		AssertStringEqual(t, testUserJohn.LastName, user.LastName)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Password, user.Password)
		testUserJohn.ID = user.ID
	})

	t.Run("Get testUserJohn By ID", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByID(testUserJohn.ID)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUserJohn.ID.String(), user.ID.String())
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Email, user.Email)
		AssertStringEqual(t, testUserJohn.Name, user.Name)
		AssertStringEqual(t, testUserJohn.LastName, user.LastName)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Password, user.Password)
	})

	t.Run("Update testUserJohn", func(t *testing.T) {
		testUserJohn.Name = "UpdatedName"
		testUserJohn.LastName = "UpdatedLastName"
		testUserJohn.Username = "updatedusername"
		testUserJohn.Email = "updateduser@example.com"
		err := testStorage.UserStore.UpdateUser(testUserJohn)
		AssertNoError(t, err)
	})

	t.Run("Get Updated User By ID", func(t *testing.T) {
		user, err := testStorage.UserStore.GetUserByID(testUserJohn.ID)
		AssertNoError(t, err)
		AssertNotNil(t, user)
		AssertStringEqual(t, testUserJohn.ID.String(), user.ID.String())
		AssertStringEqual(t, testUserJohn.Name, user.Name)
		AssertStringEqual(t, testUserJohn.LastName, user.LastName)
		AssertStringEqual(t, testUserJohn.Username, user.Username)
		AssertStringEqual(t, testUserJohn.Email, user.Email)
		AssertStringEqual(t, testUserJohn.Password, user.Password)
	})

	t.Run("Update testUserJohn To Initial values", func(t *testing.T) {
		testUserJohn.Name = "john"
		testUserJohn.LastName = "doe"
		testUserJohn.Username = "john_doe"
		testUserJohn.Email = "john_doe@example.com"
		testUserJohn.Password = "password"
		err := testStorage.UserStore.UpdateUser(testUserJohn)
		AssertNoError(t, err)
	})
}

func TestFollowRepo(t *testing.T) {

	// testUserJohn follows testUserAlice

	t.Run("testUserJohn Follow testUserALice", func(t *testing.T) {
		err := testStorage.FollowStore.FollowUser(testUserJohn.ID, testUserAlice.ID)
		AssertNoError(t, err)

	})

	t.Run("Get testUserAlice followers", func(t *testing.T) {
		follows, err := testStorage.FollowStore.GetFollowerByUserID(testUserAlice.ID)
		AssertNoError(t, err)
		AssertIntEqual(t, 1, len(follows))
		follow := follows[0]
		AssertStringEqual(t, testUserJohn.ID.String(), follow.UserID.String())
		AssertStringEqual(t, testUserAlice.ID.String(), follow.FollowID.String())

	})
	t.Run("Get testUserJohn Following", func(t *testing.T) {
		follows, err := testStorage.FollowStore.GetFollowingByUserID(testUserJohn.ID)
		AssertNoError(t, err)
		AssertIntEqual(t, 1, len(follows))
		follow := follows[0]
		AssertStringEqual(t, testUserJohn.ID.String(), follow.UserID.String())
		AssertStringEqual(t, testUserAlice.ID.String(), follow.FollowID.String())

	})

	t.Run("testUserJohn Unfollow testUserAlice", func(t *testing.T) {
		err := testStorage.FollowStore.UnFollowUser(testUserJohn.ID, testUserAlice.ID)
		AssertNoError(t, err)

	})

	t.Run("Get testUserAlice followers", func(t *testing.T) {
		follows, err := testStorage.FollowStore.GetFollowerByUserID(testUserAlice.ID)
		AssertNoError(t, err)
		AssertIntEqual(t, 0, len(follows))
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
