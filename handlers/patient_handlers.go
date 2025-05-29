package handlers

import (
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"

	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create patients"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		FullName        string `json:"full_name" binding:"required"`
		Address         string `json:"address"`
		Gender          string `json:"gender" binding:"required"`
		Age             int    `json:"age" binding:"required"`
		InsuranceNumber string `json:"insurance_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := factory.NewPatient(input.FullName, input.Address, input.Gender, input.Age, input.InsuranceNumber, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Preload("User").Find(patient)

	c.JSON(http.StatusCreated, patient)
}

func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db.DB.Preload("User").Find(&patients)
	c.JSON(http.StatusOK, patients)
}

func GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.Preload("User").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func UpdatePatient(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can update patients"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	var input struct {
		FullName        string `json:"full_name"`
		Address         string `json:"address"`
		Gender          string `json:"gender"`
		Age             int    `json:"age"`
		InsuranceNumber string `json:"insurance_number"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPatient, err := factory.NewPatient(
		input.FullName,
		input.Address,
		input.Gender,
		input.Age,
		input.InsuranceNumber,
		userID.(uint),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&patient).Updates(updatedPatient)
	db.DB.Model(&patient).Update("user_id", userID)

	db.DB.First(&patient, id)
	c.JSON(http.StatusOK, patient)
}

func DeletePatient(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can delete patients"})
		return
	}

	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	db.DB.Delete(&patient)
	c.JSON(http.StatusOK, gin.H{"message": "patient deleted"})
}
