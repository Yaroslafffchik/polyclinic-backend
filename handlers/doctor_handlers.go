package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
)

func CreateDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create doctors"})
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

	doctor, err := factory.NewDoctor(input.FullName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Create(doctor)
	c.JSON(http.StatusOK, doctor)
}

func GetDoctors(c *gin.Context) {
	var doctors []models.Doctor
	db.DB.Preload("Section").Find(&doctors)
	c.JSON(http.StatusOK, doctors)
}

func GetDoctor(c *gin.Context) {
	id := c.Param("id")
	var doctor models.Doctor
	if err := db.DB.Preload("Section").First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

func UpdateDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can update doctors"})
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

	updatedDoctor, err := factory.NewDoctor(input.FullName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&doctor).Updates(updatedDoctor)
	c.JSON(http.StatusOK, doctor)
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

	db.DB.Delete(&doctor)
	c.JSON(http.StatusOK, gin.H{"message": "doctor deleted"})
}
