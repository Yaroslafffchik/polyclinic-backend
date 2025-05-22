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

	patient, err := factory.NewPatient(input.FullName, input.Address, input.Gender, input.Age, input.InsuranceNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(patient)
	c.JSON(http.StatusOK, patient)
}

func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db.DB.Find(&patients)
	c.JSON(http.StatusOK, patients)
}

func GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.First(&patient, id).Error; err != nil {
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

	updatedPatient, err := factory.NewPatient(input.FullName, input.Address, input.Gender, input.Age, input.InsuranceNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&patient).Updates(updatedPatient)
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
