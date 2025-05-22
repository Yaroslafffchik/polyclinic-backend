package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("Role") // Для упрощения берем роль из заголовка
		if role != "doctor" && role != "registrar" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}
		c.Set("role", role)
		c.Next()
	}
}
