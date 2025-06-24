package middleware

import (
	"fmt"
	"net/http"

	"github.com/Blaze5333/ecommerce_go/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "No authorization token provided"})
			c.Abort()
			return
		}
		token = token[len("bearer "):]
		claims, msg := tokens.ValidateToken(token)
		if msg != "" {
			fmt.Println("Error validating token:", msg)
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusUnauthorized, gin.H{"message": msg})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_id", claims.Uid)
		c.Next()
	}
}
