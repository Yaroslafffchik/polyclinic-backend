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
	"strings"
	"time"
)

func CreateDoctor(c *gin.Context) {
	role := c.GetString("role")
	if role != "registrar" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only registrars can create doctors"})
		return
	}

	var input struct {
		LastName       string `json:"last_name" binding:"required"`
		FirstName      string `json:"first_name" binding:"required"`
		MiddleName     string `json:"middle_name"`
		Category       string `json:"category" binding:"required"`
		BirthDate      string `json:"birth_date" binding:"required"`
		Specialization string `json:"specialization" binding:"required"`
		Experience     int    `json:"experience" binding:"required"`
		SectionIDs     []uint `json:"section_ids" binding:"required"`
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

	doctor, err := factory.NewDoctor(input.LastName, input.FirstName, input.MiddleName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionIDs, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Привязываем участки
	for _, sectionID := range input.SectionIDs {
		db.DB.Create(&models.DoctorSections{DoctorID: doctor.ID, SectionID: sectionID})
	}

	db.DB.Preload("Sections").Preload("User").First(doctor, doctor.ID)
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
		LastName       string `json:"last_name" binding:"required"`
		FirstName      string `json:"first_name" binding:"required"`
		MiddleName     string `json:"middle_name"`
		Category       string `json:"category" binding:"required"`
		BirthDate      string `json:"birth_date" binding:"required"`
		Specialization string `json:"specialization" binding:"required"`
		Experience     int    `json:"experience" binding:"required"`
		SectionIDs     []uint `json:"section_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDoctor, err := factory.NewDoctor(input.LastName, input.FirstName, input.MiddleName, input.Category, input.BirthDate, input.Specialization, input.Experience, input.SectionIDs, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем основные поля врача
	if err := db.DB.Model(&doctor).Updates(updatedDoctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем участки
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Delete(&models.DoctorSections{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sections"})
		return
	}
	for _, sectionID := range input.SectionIDs {
		if err := db.DB.Create(&models.DoctorSections{DoctorID: doctor.ID, SectionID: sectionID}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sections"})
			return
		}
	}

	db.DB.Preload("Sections").Preload("User").First(&doctor, id)
	c.JSON(http.StatusOK, doctor)
}

func GetDoctors(c *gin.Context) {
	var doctors []models.Doctor
	if err := db.DB.Preload("Sections").Preload("User").Find(&doctors).Error; err != nil {
		log.Printf("Error loading doctors: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load doctors"})
		return
	}

	type DoctorResponse struct {
		models.Doctor
		FullName string `json:"full_name"`
	}

	var response []DoctorResponse
	for _, doctor := range doctors {
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s %s", doctor.LastName, doctor.FirstName, doctor.MiddleName))
		response = append(response, DoctorResponse{
			Doctor:   doctor,
			FullName: fullName,
		})
	}

	c.JSON(http.StatusOK, response)
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
	if err := db.DB.Preload("Sections").Preload("User").First(&doctor, doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}

	var schedules []models.Schedule
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Preload("Section").Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load schedules"})
		return
	}

	var patients []models.Patient
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Preload("Doctor").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load patients"})
		return
	}

	// Подсчет количества посещений за последний месяц
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)
	var visitCount int64
	// Поле date в модели Visit имеет тип string и формат varchar(10), предполагаем YYYY-MM-DD
	if err := db.DB.Model(&models.Visit{}).
		Where("doctor_id = ? AND date >= ? AND date < ?", doctor.ID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Count(&visitCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count visits"})
		return
	}

	// Формируем DoctorResponse с полем full_name и visit_count
	type DoctorResponse struct {
		models.Doctor
		FullName   string `json:"full_name"`
		VisitCount int64  `json:"visit_count"`
	}

	doctorResponse := DoctorResponse{
		Doctor:     doctor,
		FullName:   strings.TrimSpace(fmt.Sprintf("%s %s %s", doctor.LastName, doctor.FirstName, doctor.MiddleName)),
		VisitCount: visitCount,
	}

	// Формируем PatientResponse с полем full_name
	type PatientResponse struct {
		models.Patient
		FullName string `json:"full_name"`
	}

	var patientsResponse []PatientResponse
	for _, patient := range patients {
		patientsResponse = append(patientsResponse, PatientResponse{
			Patient:  patient,
			FullName: strings.TrimSpace(fmt.Sprintf("%s %s %s", patient.LastName, patient.FirstName, patient.MiddleName)),
		})
	}

	// Формируем ScheduleResponse с полями days, time, room
	type ScheduleResponse struct {
		models.Schedule
	}

	var schedulesResponse []ScheduleResponse
	for _, schedule := range schedules {
		schedulesResponse = append(schedulesResponse, ScheduleResponse{
			Schedule: schedule,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"doctor":    doctorResponse,
		"schedules": schedulesResponse,
		"patients":  patientsResponse,
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

	// Проверяем, есть ли другие врачи в каждом из участков врача
	var doctorSections []models.DoctorSections
	if err := db.DB.Where("doctor_id = ?", doctor.ID).Find(&doctorSections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load doctor sections"})
		return
	}
	for _, ds := range doctorSections {
		var remainingDoctors []models.DoctorSections
		if err := db.DB.Where("section_id = ? AND doctor_id != ?", ds.SectionID, doctor.ID).Find(&remainingDoctors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check remaining doctors"})
			return
		}
		if len(remainingDoctors) == 0 {
			log.Printf("Warning: Doctor ID %d was the last in Section ID %d at %s. Administrator notification required.", doctor.ID, ds.SectionID, time.Now().Format(time.RFC3339))
		}
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
	if err := db.DB.Where("doctor_id IN ?", doctorIDs).Preload("Doctor").Preload("Section").Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load schedules"})
		return
	}

	type ScheduleResponse struct {
		models.Schedule
		DoctorFullName string `json:"doctor_name"`
	}

	var response []ScheduleResponse
	for _, schedule := range schedules {
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s %s", schedule.Doctor.LastName, schedule.Doctor.FirstName, schedule.Doctor.MiddleName))
		response = append(response, ScheduleResponse{
			Schedule:       schedule,
			DoctorFullName: fullName,
		})
	}

	c.JSON(http.StatusOK, response)
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
		doctorName := strings.TrimSpace(fmt.Sprintf("%s %s %s", doctor.LastName, doctor.FirstName, doctor.MiddleName))
		stats[doctor.ID] = struct {
			DoctorName string `json:"doctor_name"`
			VisitCount int64  `json:"visit_count"`
		}{
			DoctorName: doctorName,
			VisitCount: count,
		}
	}

	c.JSON(http.StatusOK, stats)
}
