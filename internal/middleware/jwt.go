package middleware

import (
	"net/http"

	"cmsApp/pkg/jwt"

	"cmsApp/internal/constant"
	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": constant.TOKEN_NIL,
				"status":  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}
		payload, err := jwt.Check(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
		}
		c.Set("uid", payload.Id)
		c.Next()
	}
}
