package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
)

func NewPatient(fullName, address, gender string, age int, insuranceNumber string, userID uint) (*models.Patient, error) {
	// Проверка ФИО
	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}
	// Проверка адреса
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}
	// Проверка пола
	if gender != "M" && gender != "F" {
		return nil, errors.New("gender must be 'M' or 'F'")
	}
	// Проверка возраста
	if age < 0 || age > 120 {
		return nil, errors.New("age must be between 0 and 120")
	}
	// Проверка номера полиса
	if !regexp.MustCompile(`^\d{16}$`).MatchString(insuranceNumber) {
		return nil, errors.New("insurance number must be 16 digits")
	}
	// Проверка ID пользователя
	if userID == 0 {
		return nil, errors.New("user ID must be provided")
	}

	return &models.Patient{
		FullName:        fullName,
		Address:         address,
		Gender:          gender,
		Age:             age,
		InsuranceNumber: insuranceNumber,
		UserID:          userID,
	}, nil
}
