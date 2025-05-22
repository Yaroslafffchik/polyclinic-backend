package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
)

func CreateVisit(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" && role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars or doctors can create visits"})
		return
	}

	var input struct {
		PatientID    uint   `json:"patient_id"`
		DoctorID     uint   `json:"doctor_id"`
		Date         string `json:"date"`
		Complaints   string `json:"complaints"`
		Diagnosis    string `json:"diagnosis"`
		Prescription string `json:"prescription"`
		SickLeave    bool   `json:"sick_leave"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	visit, err := factory.NewVisit(input.PatientID, input.DoctorID, input.Date, input.Complaints, input.Diagnosis, input.Prescription, input.SickLeave)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(visit)
	c.JSON(http.StatusOK, visit)
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
		PatientID    uint   `json:"patient_id"`
		DoctorID     uint   `json:"doctor_id"`
		Date         string `json:"date"`
		Complaints   string `json:"complaints"`
		Diagnosis    string `json:"diagnosis"`
		Prescription string `json:"prescription"`
		SickLeave    bool   `json:"sick_leave"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedVisit, err := factory.NewVisit(input.PatientID, input.DoctorID, input.Date, input.Complaints, input.Diagnosis, input.Prescription, input.SickLeave)
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
