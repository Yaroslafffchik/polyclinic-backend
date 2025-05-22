package tests

import (
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"testing"
)

func TestCreateDoctor(t *testing.T) {
	db.Init()
	doctor, err := factory.NewDoctor("Jane Doe", "Senior", "1990-01-01", "Cardiology", 10, 1)
	if err != nil {
		t.Fatalf("Failed to create doctor: %v", err)
	}
	db.DB.Create(doctor)

	var retrieved models.Doctor
	if err := db.DB.Preload("Section").First(&retrieved, doctor.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve doctor: %v", err)
	}
	if retrieved.FullName != "Jane Doe" {
		t.Errorf("Expected full name Jane Doe, got %s", retrieved.FullName)
	}
}
