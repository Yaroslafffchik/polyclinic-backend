package factory

import (
	"errors"
	"polyclinic-backend/models"
)

func NewSection(name, address string) (*models.Section, error) {
	// Проверка названия
	if name == "" {
		return nil, errors.New("section name cannot be empty")
	}
	// Проверка адреса
	if address == "" {
		return nil, errors.New("section address cannot be empty")
	}

	return &models.Section{
		Name:    name,
		Address: address,
	}, nil
}
