package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		fields := []zap.Field{
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(ctx.Errors) > 0 {
			for _, err := range ctx.Errors {
				logger.Error("request failed", append(fields, zap.Error(err))...)
			}
		} else {
			logger.Info(path, fields...)
		}

	}
}
