package middleware

import (
	"github.com/gin-gonic/gin"
	"my-blog/pkg/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未携带 token"})
			c.Abort()
			return
		}

		// Bearer token_value
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token 格式错误"})
			c.Abort()
			return
		}

		// 解析 token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 保存用户ID到上下文中
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
