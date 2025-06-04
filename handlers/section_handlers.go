package handlers

import (
	"net/http"
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateSection(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create sections"})
		return
	}

	var input struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section, err := factory.NewSection(input.Name, input.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, section)
}

func GetSections(c *gin.Context) {
	var sections []models.Section
	db.DB.Find(&sections)
	c.JSON(http.StatusOK, sections)
}

func GetSection(c *gin.Context) {
	id := c.Param("id")
	var section models.Section
	if err := db.DB.First(&section, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "section not found"})
		return
	}
	c.JSON(http.StatusOK, section)
}

func DeleteSection(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can delete sections"})
		return
	}

	id := c.Param("id")
	var section models.Section
	if err := db.DB.First(&section, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "section not found"})
		return
	}

	var doctorCount int64
	db.DB.Model(&models.Doctor{}).Where("section_id = ?", id).Count(&doctorCount)
	if doctorCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete section with associated doctors"})
		return
	}

	if err := db.DB.Delete(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
