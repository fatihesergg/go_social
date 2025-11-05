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
	replyController *controller.ReplyController) {
	base := engine.Group("/api/v1")

	base.POST("/signup", userController.Signup)
	base.POST("/login", userController.Login)

	userRouter := base.Group("/users")
	userRouter.Use(middleware.AuthMiddleware())
	userRouter.GET("/:id", userController.GetUserByID)
	userRouter.GET("/:id/posts", userController.GetUsersPosts)
	userRouter.GET("/me", userController.GetMe)
	userRouter.POST("/:id/follow", userController.FollowUser)
	userRouter.DELETE("/:id/unfollow", userController.UnfollowUser)
	userRouter.GET("/:id/followers", userController.GetFollowerByUserID)
	userRouter.GET("/:id/following", userController.GetFollowingByUserID)
	userRouter.POST("/reset_password", userController.ResetPassword)
	userRouter.GET("/search/:username", userController.SearchUserByUsername)

	postRouter := base.Group("/posts")
	postRouter.Use(middleware.AuthMiddleware())

	postRouter.GET("/:id", postController.GetPostByID)
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
}
