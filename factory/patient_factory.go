package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
)

func NewPatient(fullName, address, gender string, age int, insuranceNumber string) (*models.Patient, error) {
	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}
	if !regexp.MustCompile(`^\d{16}$`).MatchString(insuranceNumber) {
		return nil, errors.New("insurance number must be 16 digits")
	}
	return &models.Patient{
		FullName:        fullName,
		Address:         address,
		Gender:          gender,
		Age:             age,
		InsuranceNumber: insuranceNumber,
	}, nil
}
