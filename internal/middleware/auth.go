package middleware

import (
	"net/http"

	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || len(token) < 7 {
			c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: "unauthorized"})
			c.Abort()
			return
		}
		token = token[7:] // Remove "Bearer " prefix
		claims, err := util.ParseJWT(token)
		if err != nil {

			c.JSON(http.StatusUnauthorized, util.ErrorResponse{Error: "unauthorized"})
			c.Abort()
			return
		}
		userID, _ := uuid.Parse(claims.Subject)
		c.Set("userID", userID)
		c.Next()
	}
}
