package factory

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"polyclinic-backend/models"
)

func NewUser(username, password, role string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if role != "doctor" && role != "registrar" {
		return nil, errors.New("role must be 'doctor' or 'registrar'")
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}, nil
}
