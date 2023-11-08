package middlewares

import (
	"fmt"
	"net/http"

	"bognar.dev-backend/utils"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Hello from Middleware")
		err := token.ValidateToken(c)
		if err != nil {
			fmt.Println("Error in Validate")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
			return
		}
		fmt.Println("Passed middleware")
		c.Next()

	}
}
