package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"strconv"
)

func CreatePatient(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create patients"})
		return
	}

	var input struct {
		LastName        string `json:"last_name" binding:"required"`
		FirstName       string `json:"first_name" binding:"required"`
		MiddleName      string `json:"middle_name"`
		Address         string `json:"address" binding:"required"`
		Gender          string `json:"gender" binding:"required"`
		Age             int    `json:"age" binding:"required"`
		InsuranceNumber string `json:"insurance_number" binding:"required"`
		DoctorID        uint   `json:"doctor_id"` // Теперь обычное поле
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doctorID uint = input.DoctorID
	if doctorID != 0 {
		var doctor models.Doctor
		if err := db.DB.First(&doctor, doctorID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "doctor not found"})
			return
		}
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	patient, err := factory.NewPatient(input.LastName, input.FirstName, input.MiddleName, input.Address, input.Gender, input.Age, input.InsuranceNumber, doctorID, userID.(uint))
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
	if role != "registrar" && role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars or doctors can update patients"})
		return
	}

	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	var input struct {
		LastName        string      `json:"last_name" binding:"required"`
		FirstName       string      `json:"first_name" binding:"required"`
		MiddleName      string      `json:"middle_name"`
		Address         string      `json:"address" binding:"required"`
		Gender          string      `json:"gender" binding:"required"`
		Age             interface{} `json:"age" binding:"required"`
		InsuranceNumber string      `json:"insurance_number" binding:"required"`
		DoctorID        interface{} `json:"doctor_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразование age
	var age int
	switch v := input.Age.(type) {
	case float64:
		age = int(v)
	case string:
		parsedAge, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age format"})
			return
		}
		age = parsedAge
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age type"})
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

	// Обновляем пациента
	db.DB.Model(&patient).Updates(map[string]interface{}{
		"last_name":        input.LastName,
		"first_name":       input.FirstName,
		"middle_name":      input.MiddleName,
		"address":          input.Address,
		"gender":           input.Gender,
		"age":              age,
		"insurance_number": input.InsuranceNumber,
		"doctor_id":        doctorID,
	})

	db.DB.Preload("Doctor").First(&patient, patient.ID)
	c.JSON(http.StatusOK, patient)
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

func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db.DB.Preload("User").Preload("Doctor").Find(&patients)
	c.JSON(http.StatusOK, patients)
}
