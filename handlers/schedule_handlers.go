package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
)

func CreateSchedule(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create schedules"})
		return
	}

	var input struct {
		DoctorID  uint   `json:"doctor_id" binding:"required"`
		SectionID uint   `json:"section_id" binding:"required"`
		Days      string `json:"days" binding:"required"`
		Time      string `json:"time" binding:"required"`
		Room      string `json:"room" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule, err := factory.NewSchedule(input.DoctorID, input.SectionID, input.Days, input.Time, input.Room)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Preload("Doctor").Preload("Section").First(schedule, schedule.ID)
	c.JSON(http.StatusCreated, schedule)
}

func GetSchedules(c *gin.Context) {
	var schedules []models.Schedule
	db.DB.Preload("Doctor").Find(&schedules)
	c.JSON(http.StatusOK, schedules)
}

func GetSchedule(c *gin.Context) {
	id := c.Param("id")
	var schedule models.Schedule
	if err := db.DB.Preload("Doctor").First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func UpdateSchedule(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can update schedules"})
		return
	}

	id := c.Param("id")
	var schedule models.Schedule
	if err := db.DB.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}

	var input struct {
		DoctorID  uint   `json:"doctor_id" binding:"required"`
		SectionID uint   `json:"section_id" binding:"required"`
		Days      string `json:"days" binding:"required"`
		Time      string `json:"time" binding:"required"`
		Room      string `json:"room" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSchedule, err := factory.NewSchedule(input.DoctorID, input.SectionID, input.Days, input.Time, input.Room)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&schedule).Updates(updatedSchedule)
	c.JSON(http.StatusOK, schedule)
}

func DeleteSchedule(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can delete schedules"})
		return
	}

	id := c.Param("id")
	var schedule models.Schedule
	if err := db.DB.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}

	db.DB.Delete(&schedule)
	c.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}
