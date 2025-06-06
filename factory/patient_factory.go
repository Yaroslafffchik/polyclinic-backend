package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
	"time"
)

func NewPatient(lastName, firstName, middleName, address, gender string, age int, insuranceNumber string, doctorID, userID uint) (*models.Patient, error) {
	if lastName == "" {
		return nil, errors.New("last name is required")
	}
	if firstName == "" {
		return nil, errors.New("first name is required")
	}
	if address == "" {
		return nil, errors.New("address is required")
	}
	if gender != "M" && gender != "F" {
		return nil, errors.New("gender must be M or F")
	}
	if age <= 0 {
		return nil, errors.New("age must be greater than 0")
	}
	if !regexp.MustCompile(`^\d{16}$`).MatchString(insuranceNumber) {
		return nil, errors.New("insurance number must be 16 digits")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	return &models.Patient{
		LastName:        lastName,
		FirstName:       firstName,
		MiddleName:      middleName,
		Address:         address,
		Gender:          gender,
		Age:             age,
		InsuranceNumber: insuranceNumber,
		CardCreatedAt:   time.Now(),
		UserID:          userID,
		DoctorID:        doctorID,
	}, nil
}
