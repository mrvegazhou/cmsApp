package middleware

import (
	"cmsApp/internal/constant"
	"cmsApp/pkg/rateLimit"
	"cmsApp/pkg/redisClient"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("uid")
		if userId == "" || !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": constant.USER_NOT_EXISTS,
				"status":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		client := redisClient.GetRedisClient()
		limiter := rateLimit.NewLimiter(client)
		res, err := limiter.Allow(context.Background(), fmt.Sprintf("r:%s", userId), rateLimit.PerMinute(20))
		if res.Allowed == 0 || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}
