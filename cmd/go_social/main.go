package main

import (
	"database/sql"
	"fmt"
	"log"

	docs "github.com/fatihesergg/go_social/docs"
	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/middleware"
	"github.com/fatihesergg/go_social/internal/routes"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type App struct {
	Router *gin.Engine
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
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Swagger info
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Title = "Go Social API"
	docs.SwaggerInfo.Description = "This is a simple social media API built with Go and Gin."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.DisableCaller = true
	loggerCfg.DisableStacktrace = true
	logger, err := loggerCfg.Build()

	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	engine.Use(middleware.LoggerMiddleware(logger.Named("logger_middleware")))
	engine.Use(gin.Recovery())

	err = godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	util.ApiConfig = util.LoadConfig()
	util.ApiConfig.Validate()

	DB_URI := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", util.ApiConfig.PGUser, util.ApiConfig.PGPass, util.ApiConfig.PGDB)
	db, err := sql.Open("postgres", DB_URI)
	if err != nil {
		logger.Fatal("Invalid postgres arguments")
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("Error connecting database")
	}

	userStore := database.NewUserStore(db, logger)
	postStore := database.NewPostStore(db, logger)
	commentStore := database.NewCommentStore(db, logger)
	followStore := database.NewFollowStore(db, logger)
	feedStore := database.NewFeedStore(db, logger)
	likeStore := database.NewLikeStore(db, logger)
	replyStore := database.NewReplyStore(db, logger)

	storage := database.NewPostgresStorage(userStore, postStore, commentStore, followStore, feedStore, likeStore, replyStore)
	userService := services.NewUserService(storage, logger)
	postService := services.NewPostService(storage, logger)
	commentService := services.NewCommentService(storage, logger)
	feedService := services.NewFeedService(storage, logger)
	likeService := services.NewLikeService(storage, logger)
	replyService := services.NewReplyService(storage, logger)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	rateLimiter := middleware.NewRateLimiter(1, 10)
	engine.Use(rateLimiter.TokenBucketMiddleware())

	userController := controller.NewUserController(userService)
	postController := controller.NewPostController(postService)
	commentController := controller.NewCommentController(commentService)
	feedController := controller.NewFeedController(feedService)
	likeController := controller.NewLikeController(likeService)
	replyController := controller.NewReplyController(replyService)

	routes.MountRoutes(engine, userController, postController, commentController, likeController, feedController, replyController)

	if err := engine.Run(":3000"); err != nil {
		logger.Fatal("Error starting the server")
	}

}
