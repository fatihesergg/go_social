package routes

import (
	"github.com/fatihesergg/go_social/internal/controller"
	"github.com/fatihesergg/go_social/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MountRoutes(engine *gin.Engine,
	userController *controller.UserController,
	postController *controller.PostController,
	commentController *controller.CommentController,
	likeController *controller.LikeController,
	feedController *controller.FeedController,
	replyController *controller.ReplyController,
	followController *controller.FollowController,
	notificationController *controller.NotificationController) {
	base := engine.Group("/api/v1")

	base.POST("/signup", userController.Signup)
	base.POST("/login", userController.Login)

	userRouter := base.Group("/users")
	userRouter.Use(middleware.AuthMiddleware())
	userRouter.GET("/:id", userController.GetUserByID)
	userRouter.GET("/:id/posts", userController.GetUsersPosts)
	userRouter.GET("/me", userController.GetMe)
	userRouter.POST("/reset_password", userController.ResetPassword)
	userRouter.GET("/search/:username", userController.SearchUserByUsername)

	followRouter := base.Group("/follows")
	followRouter.Use(middleware.AuthMiddleware())
	followRouter.POST("/:id/follow", followController.SendFollowRequest)
	followRouter.DELETE("/:id/unfollow", followController.UnfollowUser)
	followRouter.POST("/:id/accept", followController.AcceptFollowRequest)
	followRouter.POST("/:id/reject", followController.RejectFollowRequest)
	followRouter.POST("/:id/cancel", followController.CancelFollowRequest)
	followRouter.GET("/:id/followers", followController.GetFollowerByUserID)
	followRouter.GET("/:id/following", followController.GetFollowingByUserID)

	postRouter := base.Group("/posts")
	postRouter.Use(middleware.AuthMiddleware())

	postRouter.GET("/:id", postController.GetPostByID)
	postRouter.GET("/tag/:tag", postController.GetPostsByTag)
	postRouter.GET("/", postController.GetPosts)
	postRouter.POST("/", postController.CreatePost)
	postRouter.PUT("/:id", postController.UpdatePost)
	postRouter.POST("/:id/like", likeController.LikePost)
	postRouter.DELETE("/:id/unlike", likeController.UnlikePost)
	postRouter.GET("/:id/comments", commentController.GetCommentsByPostID)

	feedRouter := base.Group("/feed")
	feedRouter.Use(middleware.AuthMiddleware())
	feedRouter.GET("/", feedController.GetFeed)

	commentRouter := base.Group("/comments")
	commentRouter.Use(middleware.AuthMiddleware())
	commentRouter.POST("/", commentController.CreateComment)
	commentRouter.PUT("/:id", commentController.UpdateComment)
	commentRouter.GET("/:id/replies", replyController.GetCommentReplies)
	commentRouter.DELETE("/:id", commentController.DeleteComment)
	commentRouter.POST("/:id/reply", replyController.ReplyComment)
	commentRouter.POST("/:id/like", likeController.LikeComment)
	commentRouter.DELETE("/:id/unlike", likeController.UnlikeComment)

	replyRouter := base.Group("/replies")
	replyRouter.Use(middleware.AuthMiddleware())
	replyRouter.GET("/:id/replies", replyController.GetRepliesByParent)
	replyRouter.PUT("/:id", replyController.UpdateReply)
	replyRouter.DELETE("/:id", replyController.DeleteReply)
	replyRouter.POST("/:id/reply", replyController.ReplyAReply)
	replyRouter.POST("/:id/like", likeController.LikeReply)
	replyRouter.DELETE("/:id/unlike", likeController.UnlikeReply)

	notificationRouter := base.Group("/notifications")
	notificationRouter.Use(middleware.AuthMiddleware())
	notificationRouter.GET("/", notificationController.GetNotification)

}
