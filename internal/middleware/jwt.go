package middleware

import (
	"cmsApp/configs"
	"net/http"
	"strings"

	"cmsApp/pkg/jwt"

	"cmsApp/internal/constant"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": constant.TOKEN_NIL,
				"status":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		payload, err := jwt.Check(token, configs.App.Login.JwtSecret, false)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}
		c.Set("uid", payload.ID)
		c.Next()
	}
}
