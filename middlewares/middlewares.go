package middlewares

import (
	"fmt"
	"net/http"

	"bognar.dev-backend/utils"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("jwt middleware")
		err := token.ValidateToken(c)
		if err != nil {
			fmt.Println("Middleware error: ", err)
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		fmt.Println("Authorized")
		c.Next()
	}
}
