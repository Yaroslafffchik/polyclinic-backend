package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"strconv"
)

func CreateVisit(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" && role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars or doctors can create visits"})
		return
	}

	var input struct {
		PatientID         interface{} `json:"patient_id" binding:"required"`
		DoctorID          interface{} `json:"doctor_id" binding:"required"`
		Date              string      `json:"date" binding:"required"`
		Complaints        string      `json:"complaints"`
		Diagnosis         string      `json:"diagnosis"`
		Prescription      string      `json:"prescription"`
		SickLeave         bool        `json:"sick_leave"`
		SickLeaveDuration interface{} `json:"sick_leave_duration"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразование patient_id
	var patientID uint
	switch v := input.PatientID.(type) {
	case float64:
		patientID = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id format"})
			return
		}
		patientID = uint(id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient_id type"})
		return
	}

	// Преобразование doctor_id
	var doctorID uint
	switch v := input.DoctorID.(type) {
	case float64:
		doctorID = uint(v)
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id format"})
			return
		}
		doctorID = uint(id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor_id type"})
		return
	}

	// Преобразование sick_leave_duration
	var sickLeaveDuration int
	switch v := input.SickLeaveDuration.(type) {
	case float64:
		sickLeaveDuration = int(v)
	case string:
		dur, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sick_leave_duration format"})
			return
		}
		sickLeaveDuration = dur
	case nil:
		sickLeaveDuration = 0 // По умолчанию 0, если поле не указано
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sick_leave_duration type"})
		return
	}

	visit, err := factory.NewVisit(patientID, doctorID, input.Date, input.Complaints, input.Diagnosis, input.Prescription, input.SickLeave, sickLeaveDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(visit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Preload("Patient").Preload("Doctor").First(visit, visit.ID)
	c.JSON(http.StatusCreated, visit)
}

func GetVisits(c *gin.Context) {
	var visits []models.Visit
	db.DB.Preload("Patient").Preload("Doctor").Find(&visits)
	c.JSON(http.StatusOK, visits)
}

func GetVisit(c *gin.Context) {
	id := c.Param("id")
	var visit models.Visit
	if err := db.DB.Preload("Patient").Preload("Doctor").First(&visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
		return
	}
	c.JSON(http.StatusOK, visit)
}

func UpdateVisit(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" && role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars or doctors can update visits"})
		return
	}

	id := c.Param("id")
	var visit models.Visit
	if err := db.DB.First(&visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
		return
	}

	var input struct {
		PatientID         uint   `json:"patient_id"`
		DoctorID          uint   `json:"doctor_id"`
		Date              string `json:"date"`
		Complaints        string `json:"complaints"`
		Diagnosis         string `json:"diagnosis"`
		Prescription      string `json:"prescription"`
		SickLeave         bool   `json:"sick_leave"`
		SickLeaveDuration int    `json:"sick_leave_duration"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedVisit, err := factory.NewVisit(input.PatientID, input.DoctorID, input.Date, input.Complaints, input.Diagnosis, input.Prescription, input.SickLeave, input.SickLeaveDuration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&visit).Updates(updatedVisit)
	c.JSON(http.StatusOK, visit)
}

func DeleteVisit(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" && role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars or doctors can delete visits"})
		return
	}

	id := c.Param("id")
	var visit models.Visit
	if err := db.DB.First(&visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "visit not found"})
		return
	}

	db.DB.Delete(&visit)
	c.JSON(http.StatusOK, gin.H{"message": "visit deleted"})
}
