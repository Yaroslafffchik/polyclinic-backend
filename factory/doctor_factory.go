package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
	"time"
)

func NewDoctor(fullName, category, birthDate, specialization string, experience int, sectionID, userID uint) (*models.Doctor, error) {
	// Проверка ФИО
	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}
	// Проверка специализации
	if specialization == "" {
		return nil, errors.New("specialization cannot be empty")
	}
	// Проверка формата даты рождения (YYYY-MM-DD)
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(birthDate) {
		return nil, errors.New("birth date must be in YYYY-MM-DD format")
	}
	// Проверка, что дата рождения в прошлом (без учета времени)
	parsedDate, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return nil, errors.New("invalid birth date format")
	}
	currentDate := time.Now().Truncate(24 * time.Hour) // Убираем время для точного сравнения дат
	if !parsedDate.Before(currentDate) {
		return nil, errors.New("birth date must be in the past")
	}

	// Проверка возраста врача (от 25 до 100 лет)
	currentYear := time.Now().Year()
	birthYear := parsedDate.Year()
	age := currentYear - birthYear
	if age < 25 || age > 100 {
		return nil, errors.New("doctor's age must be between 25 and 100 years")
	}

	// Проверка стажа
	if experience < 0 {
		return nil, errors.New("experience cannot be negative")
	}
	// Проверка, что стаж не превышает возраст минус 20 лет (примерно начало карьеры)
	if experience > (age - 20) {
		return nil, errors.New("experience cannot exceed age minus 20 years")
	}
	// Проверка ID участка
	if sectionID == 0 {
		return nil, errors.New("section ID must be provided")
	}
	// Проверка ID пользователя
	if userID == 0 {
		return nil, errors.New("user ID must be provided")
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
