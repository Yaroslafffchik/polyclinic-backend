package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
)

func CreatePatient(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create patients"})
		return
	}

	var input struct {
		FullName        string `json:"full_name" binding:"required"`
		Address         string `json:"address" binding:"required"`
		Gender          string `json:"gender" binding:"required"`
		Age             int    `json:"age" binding:"required"`
		InsuranceNumber string `json:"insurance_number" binding:"required"`
		DoctorName      string `json:"doctor_name"` // Для создания
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doctorID uint
	if input.DoctorName != "" {
		var doctor models.Doctor
		if err := db.DB.Where("full_name = ?", input.DoctorName).First(&doctor).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "doctor not found"})
			return
		}
		doctorID = doctor.ID
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	patient, err := factory.NewPatient(input.FullName, input.Address, input.Gender, input.Age, input.InsuranceNumber, doctorID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
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
		DoctorID        *uint  `json:"doctor_id"` // Используем указатель, чтобы GORM мог обработать 0
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received input: %+v\n", input) // Отладка

	// Обновляем поля
	updates := map[string]interface{}{
		"full_name":        input.FullName,
		"address":          input.Address,
		"gender":           input.Gender,
		"age":              input.Age,
		"insurance_number": input.InsuranceNumber,
		"user_id":          userID,
	}
	// Добавляем DoctorID в обновление только если он был передан
	if input.DoctorID != nil {
		updates["doctor_id"] = *input.DoctorID
	}

	if err := db.DB.Model(&patient).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Preload("Doctor").First(&patient, id)
	c.JSON(http.StatusOK, patient)
}

func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db.DB.Preload("User").Preload("Doctor").Find(&patients)
	c.JSON(http.StatusOK, patients)
}

func GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.Preload("User").Preload("Doctor").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}
	// Загружаем визиты пациента
	var visits []models.Visit
	db.DB.Where("patient_id = ?", patient.ID).Preload("Doctor").Find(&visits)
	c.JSON(http.StatusOK, gin.H{
		"patient": patient,
		"visits":  visits,
	})
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
