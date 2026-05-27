package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "github.com/fatihesergg/go_social/docs"
	"github.com/fatihesergg/go_social/internal/broker"
	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/database"
	"github.com/fatihesergg/go_social/internal/middleware"
	"github.com/fatihesergg/go_social/internal/routes"
	"github.com/fatihesergg/go_social/internal/services"
	"github.com/fatihesergg/go_social/internal/util"
	"github.com/fatihesergg/go_social/internal/worker"
	"github.com/fatihesergg/go_social/internal/ws"
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
		logger.Info("Error loading .env file, using env variables")
	}

	util.ApiConfig = util.LoadConfig()
	util.ApiConfig.Validate()

	DB_URI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", util.ApiConfig.PGUser, util.ApiConfig.PGPass, util.ApiConfig.PGHost, util.ApiConfig.PGPort, util.ApiConfig.PGDB)
	db, err := sql.Open("postgres", DB_URI)
	if err != nil {
		logger.Fatal("Invalid postgres arguments")
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("Error connecting database")
	}

	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s", util.ApiConfig.RABBITMQ_USER, util.ApiConfig.RABBITMQ_PASSWORD, util.ApiConfig.RABBITMQ_HOST, util.ApiConfig.RABBITMQ_PORT)
	rabbitmq, err := broker.NewRabbitMq(rabbitmqURL)
	fmt.Println(util.ApiConfig.RabbitMQ)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer rabbitmq.Close()

	hub := ws.NewWsHub()

	wsocket := engine.Group("/ws")
	wsocket.Use(middleware.AuthMiddleware())
	wsocket.GET("/", func(ctx *gin.Context) {
		ws.ServeWS(hub, ctx)
	})

	userStore := database.NewUserStore(db)
	postStore := database.NewPostStore(db)
	commentStore := database.NewCommentStore(db)
	followStore := database.NewFollowStore(db)
	feedStore := database.NewFeedStore(db)
	likeStore := database.NewLikeStore(db)
	replyStore := database.NewReplyStore(db)
	tagStore := database.NewTagStore(db)
	notificationStore := database.NewNotificationStore(db)

	userService := services.NewUserService(userStore, followStore, postStore)
	postService := services.NewPostService(postStore)
	commentService := services.NewCommentService(postStore, commentStore)
	feedService := services.NewFeedService(feedStore)
	likeService := services.NewLikeService(likeStore, postStore, commentStore, rabbitmq)
	replyService := services.NewReplyService(replyStore, commentStore)
	tagService := services.NewTagService(tagStore)
	followService := services.NewFollowService(followStore, userStore)
	notificationService := services.NewNotificationService(notificationStore, logger)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	rateLimiter := middleware.NewRateLimiter(1, 10)
	engine.Use(rateLimiter.TokenBucketMiddleware())

	userController := controller.NewUserController(userService)
	postController := controller.NewPostController(postService, tagService)
	commentController := controller.NewCommentController(commentService)
	feedController := controller.NewFeedController(feedService)
	likeController := controller.NewLikeController(likeService)
	replyController := controller.NewReplyController(replyService)
	followController := controller.NewFollowController(followService)
	notificationController := controller.NewNotificationController(notificationService)

	routes.MountRoutes(engine, userController, postController, commentController, likeController, feedController, replyController, followController, notificationController)

	server := &http.Server{Addr: ":3000", Handler: engine.Handler()}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting the server")
		} else {
			logger.Info("Server listening on port 3000")
		}
	}()

	_, err = rabbitmq.DeclareQueue("post_like_queue")
	if err != nil {
		logger.Error("Notification queue error", zap.Error(err))
	}

	_, err = rabbitmq.DeclareQueue("comment_like_queue")
	if err != nil {
		logger.Error("Notification queue error", zap.Error(err))
	}

	_, err = rabbitmq.DeclareQueue("reply_like_queue")
	if err != nil {
		logger.Error("Notification queue error", zap.Error(err))
	}

	err = rabbitmq.DeclareExchange("like_event", "direct", true, true)
	if err != nil {
		logger.Error("exchange declare fail", zap.Error(err))
	}

	err = rabbitmq.BindQueueToExchange("like_event", "post", "post_like_queue")
	if err != nil {
		logger.Error("binding queue to exchange fail", zap.Error(err))
	}

	err = rabbitmq.BindQueueToExchange("like_event", "comment", "comment_like_queue")
	if err != nil {
		logger.Error("binding queue to exchange fail", zap.Error(err))
	}

	err = rabbitmq.BindQueueToExchange("like_event", "reply", "reply_like_queue")
	if err != nil {
		logger.Error("binding queue to exchange fail", zap.Error(err))
	}

	notificationWorker := worker.NewNotificationWorker(rabbitmq.Channel, logger, hub, postStore, userStore, commentStore, replyStore, notificationStore)
	mainCtx, mainCancel := context.WithCancel(context.Background())
	go func() {
		err = notificationWorker.ConsumePostLike(mainCtx)
		if err != nil {
			logger.Error("error while starting to consume post like event", zap.Error(err))
		}
	}()

	go func() {
		err = notificationWorker.ConsumeCommentLike(mainCtx)
		if err != nil {
			logger.Error("error while starting to consume comment like events")
		}
	}()

	go func() {
		err = notificationWorker.ConsumeReplyLike(mainCtx)
		if err != nil {
			logger.Error("error while starting to consume reply like event", zap.Error(err))
		}
	}()

	defer mainCancel()

	go hub.Run(mainCtx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received.Shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Info("Error while shutdown server")
	}
	logger.Info("Server shutdown successfully")
}
