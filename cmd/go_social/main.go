package main

import (
	"database/sql"
	"os"

	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	engine := gin.Default()
	base := engine.Group("/api/v1")

	dotenv := godotenv.Load()
	if dotenv != nil {
		panic("Error loading .env file")
	}

	DSN := os.Getenv("DATABASE_URI")
	if DSN == "" {
		panic("DATABASE_URI is not set")
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET is not set")
	}

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		panic("Error connecting to the database")
	}

	userStore := database.NewUserStore(db)
	postStore := database.NewPostStore(db)
	storage := database.NewPostgresStorage(userStore, postStore)

	userController := controller.UserController{Storage: *storage}
	postController := controller.PostController{Storage: *storage}

	base.POST("/signup", userController.Signup)
	base.POST("/login", userController.Login)

	userRouter := base.Group("/users")
	userRouter.GET("/:id", userController.GetUserByID)

	postRouter := base.Group("/posts")
	postRouter.Use(middleware.AuthMiddleware())

	postRouter.GET("/:id", postController.GetPostByID)
	postRouter.GET("/", postController.GetPosts)
	postRouter.POST("/", postController.CreatePost)
	postRouter.PUT("/:id", postController.UpdatePost)

	if err := engine.Run(":3000"); err != nil {
		panic("Error starting the server")
	}

}
