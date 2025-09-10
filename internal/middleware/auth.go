package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fatihesergg/go_social/internal/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || len(token) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		token = token[7:] // Remove "Bearer " prefix
		claims, err := util.ParseJWT(token)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		userID, _ := strconv.Atoi(claims.Subject)
		c.Set("userID", userID)
		c.Next()
	}
}
