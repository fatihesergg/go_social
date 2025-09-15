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

	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET is not set")
	}

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		panic("Error connecting to the database")
	}

	userStore := database.NewUserStore(db)
	postStore := database.NewPostStore(db)
	commentStore := database.NewCommentStore(db)
	followStore := database.NewFollowStore(db)
	feedStore := database.NewFeedStore(db)

	storage := database.NewPostgresStorage(userStore, postStore, commentStore, followStore, feedStore)

	userController := controller.UserController{Storage: *storage}
	postController := controller.PostController{Storage: *storage}
	commentController := controller.CommentController{Storage: *storage}
	feedController := controller.FeedController{Storage: *storage}

	base.POST("/signup", userController.Signup)
	base.POST("/login", userController.Login)

	userRouter := base.Group("/users")
	userRouter.Use(middleware.AuthMiddleware())
	userRouter.GET("/:id", userController.GetUserByID)
	userRouter.GET("/:id/posts", userController.GetUsersPosts)
	userRouter.GET("/getMe", userController.GetMe)
	userRouter.POST("/:id/follow", userController.FollowUser)
	userRouter.DELETE("/:id/unfollow", userController.UnfollowUser)
	userRouter.GET("/:id/followers", userController.GetFollowerByUserID)
	userRouter.GET("/:id/following", userController.GetFollowingByUserID)

	postRouter := base.Group("/posts")
	postRouter.Use(middleware.AuthMiddleware())

	postRouter.GET("/:id", postController.GetPostByID)
	postRouter.GET("/", postController.GetPosts)
	postRouter.POST("/", postController.CreatePost)
	postRouter.PUT("/:id", postController.UpdatePost)

	feedRouter := base.Group("/feed")
	feedRouter.Use(middleware.AuthMiddleware())
	feedRouter.GET("/", feedController.GetFeed)

	commentRouter := base.Group("/comments")
	commentRouter.Use(middleware.AuthMiddleware())
	commentRouter.POST("/", commentController.CreateComment)
	commentRouter.GET("/:post_id", commentController.GetCommentsByPostID)
	commentRouter.PUT("/:id", commentController.UpdateComment)
	commentRouter.DELETE("/:id", commentController.DeleteComment)

	if err := engine.Run(":3000"); err != nil {
		panic("Error starting the server")
	}

}
