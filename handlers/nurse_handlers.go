package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
)

func CreateNurse(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create nurses"})
		return
	}

	var input struct {
		LastName   string `json:"last_name" binding:"required"`
		FirstName  string `json:"first_name" binding:"required"`
		MiddleName string `json:"middle_name"`
		SectionID  uint   `json:"section_id" binding:"required"`
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

	nurse, err := factory.NewNurse(input.LastName, input.FirstName, input.MiddleName, input.SectionID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(nurse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Preload("Section").Preload("User").First(nurse, nurse.ID)
	c.JSON(http.StatusCreated, nurse)
}

func GetNurses(c *gin.Context) {
	var nurses []models.Nurse
	db.DB.Preload("Section").Preload("User").Find(&nurses)
	c.JSON(http.StatusOK, nurses)
}

func GetNurse(c *gin.Context) {
	id := c.Param("id")
	var nurse models.Nurse
	if err := db.DB.Preload("Section").Preload("User").First(&nurse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nurse not found"})
		return
	}
	c.JSON(http.StatusOK, nurse)
}

func UpdateNurse(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can update nurses"})
		return
	}

	id := c.Param("id")
	var nurse models.Nurse
	if err := db.DB.First(&nurse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nurse not found"})
		return
	}

	var input struct {
		LastName   string `json:"last_name" binding:"required"`
		FirstName  string `json:"first_name" binding:"required"`
		MiddleName string `json:"middle_name"`
		SectionID  uint   `json:"section_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNurse, err := factory.NewNurse(input.LastName, input.FirstName, input.MiddleName, input.SectionID, nurse.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&nurse).Updates(map[string]interface{}{
		"last_name":   updatedNurse.LastName,
		"first_name":  updatedNurse.FirstName,
		"middle_name": updatedNurse.MiddleName,
		"section_id":  updatedNurse.SectionID,
	})

	db.DB.Preload("Section").Preload("User").First(&nurse, nurse.ID)
	c.JSON(http.StatusOK, nurse)
}

func DeleteNurse(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can delete nurses"})
		return
	}

	id := c.Param("id")
	var nurse models.Nurse
	if err := db.DB.First(&nurse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nurse not found"})
		return
	}

	db.DB.Delete(&nurse)
	c.JSON(http.StatusOK, gin.H{"message": "nurse deleted"})
}
