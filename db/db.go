package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"polyclinic-backend/config"
	"polyclinic-backend/models"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	DB.AutoMigrate(&models.Patient{}, &models.Doctor{}, &models.Section{}, &models.Schedule{}, &models.Visit{}, &models.User{}, &models.Log{})
}
