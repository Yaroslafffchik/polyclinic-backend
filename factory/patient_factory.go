package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
)

func NewPatient(fullName, address, gender string, age int, insuranceNumber string, doctorID, userID uint) (*models.Patient, error) {
	if fullName == "" {
		return nil, errors.New("full name is required")
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
		FullName:        fullName,
		Address:         address,
		Gender:          gender,
		Age:             age,
		InsuranceNumber: insuranceNumber,
		DoctorID:        doctorID,
		UserID:          userID,
	}, nil
}
