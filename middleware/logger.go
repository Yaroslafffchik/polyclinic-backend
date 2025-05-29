package middleware

import (
	"github.com/gin-gonic/gin"
	"polyclinic-backend/db"
	"polyclinic-backend/models"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		userID, exists := c.Get("user_id")
		if !exists {
			userID = uint(0)
		}

		logEntry := models.Log{
			UserID:   userID.(uint),
			Action:   c.Request.Method + " " + c.Request.URL.Path,
			Duration: duration.String(),
		}
		db.DB.Create(&logEntry)
	}
}
