package main

import (
	"database/sql"
	"fmt"
	"os"

	docs "github.com/fatihesergg/go_social/docs"
	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	Router  *gin.Engine
	Storage *database.Storage
}

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @title						Go Social API
// @version					1.0
// @description				This is a simple social media API built with Go and Gin.
// @host						localhost:3000
// @BasePath					/api/v1
func main() {
	engine := gin.Default()

	// Swagger info
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "Go Social API"
	docs.SwaggerInfo.Description = "This is a simple social media API built with Go and Gin."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"

	dotenv := godotenv.Load()
	if dotenv != nil {
		panic("Error loading .env file")
	}

	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	if pgUser == "" || pgPassword == "" || pgDB == "" {
		panic("Database environment variables are not set")
	}

	DSN := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", pgUser, pgPassword, pgDB)

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
	likeStore := database.NewLikeStore(db)

	storage := database.NewPostgresStorage(userStore, postStore, commentStore, followStore, feedStore, likeStore)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	rateLimiter := middleware.NewRateLimiter(1, 10)
	engine.Use(rateLimiter.TokenBucketMiddleware())
	app := App{
		Router:  engine,
		Storage: storage,
	}
	base := app.Router.Group("/api/v1")

	userController := controller.UserController{Storage: *storage}
	postController := controller.PostController{Storage: *storage}
	commentController := controller.CommentController{Storage: *storage}
	feedController := controller.FeedController{Storage: *storage}
	likeController := controller.LikeController{Storage: *storage}

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
	userRouter.POST("/reset_password", userController.ResetPassword)

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

	likeRouter := base.Group("/likes")
	likeRouter.Use(middleware.AuthMiddleware())
	likeRouter.POST("/", likeController.LikePost)
	likeRouter.DELETE("/:id", likeController.UnlikePost)

	if err := app.Router.Run(":3000"); err != nil {
		panic("Error starting the server")
	}

}
