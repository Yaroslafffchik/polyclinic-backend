package factory

import (
	"errors"
	"polyclinic-backend/models"
)

func NewDoctor(fullName, category, birthDate, specialization string, experience int, sectionID, userID uint) (*models.Doctor, error) {
	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}
	if specialization == "" {
		return nil, errors.New("specialization cannot be empty")
	}
	return &models.Doctor{
		FullName:       fullName,
		Category:       category,
		Experience:     experience,
		BirthDate:      birthDate,
		Specialization: specialization,
		SectionID:      sectionID,
		UserID:         userID,
	}, nil
}
