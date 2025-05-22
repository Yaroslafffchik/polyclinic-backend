package middleware

import (
	"polyclinic-backend/db"
	"polyclinic-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		logEntry := models.Log{
			UserID:    1, // Для MVP захардкодим user_id
			Action:    c.Request.Method + " " + c.Request.URL.Path,
			Timestamp: time.Now().Format(time.RFC3339),
			Duration:  duration.String(), // Добавляем длительность
		}
		db.DB.Create(&logEntry)
	}
}
