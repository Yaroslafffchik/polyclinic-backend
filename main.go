package main

import (
	"polyclinic-backend/db"
	"polyclinic-backend/handlers"
	"polyclinic-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()
	r.Use(middleware.Logger())
	r.Use(middleware.Auth())

	// Patients
	r.GET("/patients", handlers.GetPatients)
	r.GET("/patients/:id", handlers.GetPatient)
	r.POST("/patients", handlers.CreatePatient)
	r.PUT("/patients/:id", handlers.UpdatePatient)
	r.DELETE("/patients/:id", handlers.DeletePatient)

	// Doctors
	r.GET("/doctors", handlers.GetDoctors)
	r.GET("/doctors/:id", handlers.GetDoctor)
	r.POST("/doctors", handlers.CreateDoctor)
	r.PUT("/doctors/:id", handlers.UpdateDoctor)
	r.DELETE("/doctors/:id", handlers.DeleteDoctor)

	// Schedules
	r.GET("/schedules", handlers.GetSchedules)
	r.GET("/schedules/:id", handlers.GetSchedule)
	r.POST("/schedules", handlers.CreateSchedule)
	r.PUT("/schedules/:id", handlers.UpdateSchedule)
	r.DELETE("/schedules/:id", handlers.DeleteSchedule)

	// Visits
	r.GET("/visits", handlers.GetVisits)
	r.GET("/visits/:id", handlers.GetVisit)
	r.POST("/visits", handlers.CreateVisit)
	r.PUT("/visits/:id", handlers.UpdateVisit)
	r.DELETE("/visits/:id", handlers.DeleteVisit)

	r.Run(":8080")
}
