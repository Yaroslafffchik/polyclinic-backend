package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
	"time"
)

func NewDoctor(lastName, firstName, middleName, category, birthDate, specialization string, experience int, sectionIDs []uint, userID uint) (*models.Doctor, error) {
	if lastName == "" || firstName == "" {
		return nil, errors.New("last name and first name cannot be empty")
	}
	if specialization == "" {
		return nil, errors.New("specialization cannot be empty")
	}
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(birthDate) {
		return nil, errors.New("birth date must be in YYYY-MM-DD format")
	}
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return nil, errors.New("invalid birth date format")
	}
	currentDate := time.Now().Truncate(24 * time.Hour)
	if !parsedDate.Before(currentDate) {
		return nil, errors.New("birth date must be in the past")
	}
	currentYear := time.Now().Year()
	birthYear := parsedDate.Year()
	age := currentYear - birthYear
	if age < 25 || age > 100 {
		return nil, errors.New("doctor's age must be between 25 and 100 years")
	}
	if experience < 0 {
		return nil, errors.New("experience cannot be negative")
	}
	if experience > (age - 20) {
		return nil, errors.New("experience cannot exceed age minus 20 years")
	}
	if len(sectionIDs) == 0 {
		return nil, errors.New("at least one section ID must be provided")
	}
	if userID == 0 {
		return nil, errors.New("user ID must be provided")
	}

	return &models.Doctor{
		LastName:       lastName,
		FirstName:      firstName,
		MiddleName:     middleName,
		Category:       category,
		Experience:     experience,
		BirthDate:      birthDate,
		Specialization: specialization,
		UserID:         userID,
	}, nil
}
