package factory

import (
	"errors"
	"polyclinic-backend/models"
)

func NewNurse(lastName, firstName, middleName string, sectionID, userID uint) (*models.Nurse, error) {
	if lastName == "" || firstName == "" {
		return nil, errors.New("last name and first name cannot be empty")
	}
	if sectionID == 0 {
		return nil, errors.New("section ID must be provided")
	}
	if userID == 0 {
		return nil, errors.New("user ID must be provided")
	}

	return &models.Nurse{
		LastName:   lastName,
		FirstName:  firstName,
		MiddleName: middleName,
		SectionID:  sectionID,
		UserID:     userID,
	}, nil
}
