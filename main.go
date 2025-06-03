package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"polyclinic-backend/db"
	"polyclinic-backend/handlers"
	"polyclinic-backend/middleware"
)

func main() {
	db.Init()

	r := gin.Default()
	r.Use(middleware.Logger())

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},                   // Разрешаем фронтенд
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Добавляем OPTIONS
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 часов
	}))

	// Открытые маршруты
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Защищённые маршруты
	api := r.Group("/api")
	api.Use(middleware.Auth())
	{
		// Patients
		api.GET("/patients", handlers.GetPatients)
		api.GET("/patients/:id", handlers.GetPatient)
		api.POST("/patients", handlers.CreatePatient)
		api.PUT("/patients/:id", handlers.UpdatePatient)
		api.DELETE("/patients/:id", handlers.DeletePatient)

		// Doctors
		api.GET("/doctors", handlers.GetDoctors)
		api.GET("/doctors/:id", handlers.GetDoctor)
		api.POST("/doctors", handlers.CreateDoctor)
		api.PUT("/doctors/:id", handlers.UpdateDoctor)
		api.DELETE("/doctors/:id", handlers.DeleteDoctor)

		// Schedules
		api.GET("/schedules", handlers.GetSchedules)
		api.GET("/schedules/:id", handlers.GetSchedule)
		api.POST("/schedules", handlers.CreateSchedule)
		api.PUT("/schedules/:id", handlers.UpdateSchedule)
		api.DELETE("/schedules/:id", handlers.DeleteSchedule)

		// Visits
		api.GET("/visits", handlers.GetVisits)
		api.GET("/visits/:id", handlers.GetVisit)
		api.POST("/visits", handlers.CreateVisit)
		api.PUT("/visits/:id", handlers.UpdateVisit)
		api.DELETE("/visits/:id", handlers.DeleteVisit)

		// Sections
		api.GET("/sections", handlers.GetSections)
		api.GET("/sections/:id", handlers.GetSection)
		api.POST("/sections", handlers.CreateSection)
		api.DELETE("/sections/:id", handlers.DeleteSection)
	}

	r.Run(":8080")
}
