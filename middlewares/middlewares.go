package middlewares

import (
	"github.com/juju/ratelimit"
	"net/http"

	"bognar.dev-backend/utils"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

var limiter = ratelimit.NewBucketWithRate(50, 100)

func RateLimit(c *gin.Context) {
	if limiter.TakeAvailable(1) == 0 {

		c.JSON(http.StatusTooManyRequests, gin.H{"status": "Too many requests"})
	}
}
