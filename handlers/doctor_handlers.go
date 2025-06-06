package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"strconv"
	"time"
)

func CreateDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create doctors"})
		return
	}

	var input struct {
		FullName       string `json:"full_name" binding:"required"`
		Category       string `json:"category" binding:"required"`
		BirthDate      string `json:"birth_date" binding:"required"`
		Specialization string `json:"specialization" binding:"required"`
		Experience     int    `json:"experience" binding:"required"`
		SectionID      uint   `json:"section_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	doctor, err := factory.NewDoctor(input.FullName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Preload("Section").Preload("User").First(doctor, doctor.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load related data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

func UpdateDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can update doctors"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var doctor models.Doctor
	if err := db.DB.First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}

	var input struct {
		FullName       string `json:"full_name"`
		Category       string `json:"category"`
		Experience     int    `json:"experience"`
		BirthDate      string `json:"birth_date"`
		Specialization string `json:"specialization"`
		SectionID      uint   `json:"section_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDoctor, err := factory.NewDoctor(input.FullName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&doctor).Updates(updatedDoctor)
	db.DB.Model(&doctor).Update("user_id", userID)

	db.DB.Preload("Section").First(&doctor, id)
	c.JSON(http.StatusOK, doctor)
}

func GetDoctors(c *gin.Context) {
	var doctors []models.Doctor
	db.DB.Preload("Section").Preload("User").Find(&doctors)
	c.JSON(http.StatusOK, doctors)
}

func GetDoctor(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID format"})
		return
	}

	var doctor models.Doctor
	if err := db.DB.Preload("Section").Preload("User").First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	var schedules []models.Schedule
	db.DB.Where("doctor_id = ?", doctor.ID).Find(&schedules)
	var patients []models.Patient
	db.DB.Where("doctor_id = ?", doctor.ID).Preload("Doctor").Find(&patients)
	c.JSON(http.StatusOK, gin.H{
		"doctor":    doctor,
		"schedules": schedules,
		"patients":  patients,
	})
}

func DeleteDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can delete doctors"})
		return
	}

	id := c.Param("id")
	var doctor models.Doctor
	if err := db.DB.First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}

	var replacementDoctor models.Doctor
	if err := db.DB.Where("specialization = ? AND id != ?", doctor.Specialization, id).
		Order("experience desc").First(&replacementDoctor).Error; err != nil {
		if err := db.DB.Order("experience desc").First(&replacementDoctor).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no available doctor for reassignment"})
			return
		}
	}

	var patients []models.Patient
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load patients"})
		return
	}
	for _, patient := range patients {
		patient.DoctorID = replacementDoctor.ID
		if err := db.DB.Save(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reassign patient"})
			return
		}
	}

	var remainingDoctors []models.Doctor
	db.DB.Where("section_id = ? AND id != ?", doctor.SectionID, id).Find(&remainingDoctors)
	if len(remainingDoctors) == 0 {
		log.Printf("Warning: Doctor ID %d was the last in Section ID %d at %s. Administrator notification required.", doctor.ID, doctor.SectionID, time.Now().Format(time.RFC3339))
	}

	if err := db.DB.Delete(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("doctor deleted, patients reassigned to doctor ID %d", replacementDoctor.ID)})
}

func GetDoctorPatients(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID format"})
		return
	}

	var doctor models.Doctor
	if err := db.DB.First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	var patients []models.Patient
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Preload("Doctor").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load patients"})
		return
	}
	c.JSON(http.StatusOK, patients)
}

func GetSchedulesBySpecialization(c *gin.Context) {
	specialization := c.Param("specialization")
	var doctors []models.Doctor
	if err := db.DB.Where("specialization = ?", specialization).Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load doctors"})
		return
	}
	if len(doctors) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no doctors found with this specialization"})
		return
	}

	var doctorIDs []uint
	for _, doctor := range doctors {
		doctorIDs = append(doctorIDs, doctor.ID)
	}

	var schedules []models.Schedule
	if err := db.DB.Where("doctor_id IN ?", doctorIDs).Preload("Doctor").Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func GetDoctorVisitStats(c *gin.Context) {
	endDate := time.Date(2025, time.June, 1, 0, 0, 0, 0, time.UTC)
	startDate := endDate.AddDate(0, -1, 0)

	var doctors []models.Doctor
	if err := db.DB.Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load doctors"})
		return
	}

	stats := make(map[uint]struct {
		DoctorName string `json:"doctor_name"`
		VisitCount int64  `json:"visit_count"`
	})

	for _, doctor := range doctors {
		var count int64
		if err := db.DB.Model(&models.Visit{}).
			Where("doctor_id = ? AND date >= ? AND date < ?", doctor.ID, startDate, endDate).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count visits"})
			return
		}
		stats[doctor.ID] = struct {
			DoctorName string `json:"doctor_name"`
			VisitCount int64  `json:"visit_count"`
		}{
			DoctorName: doctor.FullName,
			VisitCount: count,
		}
	}

	c.JSON(http.StatusOK, stats)
}
