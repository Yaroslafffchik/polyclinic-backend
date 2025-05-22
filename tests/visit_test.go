package tests

import (
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"testing"
)

func TestCreateVisit(t *testing.T) {
	db.Init()
	visit, err := factory.NewVisit(1, 1, "2025-05-22", "Fever", "Flu", "Rest", true)
	if err != nil {
		t.Fatalf("Failed to create visit: %v", err)
	}
	db.DB.Create(visit)

	var retrieved models.Visit
	if err := db.DB.Preload("Patient").Preload("Doctor").First(&retrieved, visit.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve visit: %v", err)
	}
	if retrieved.Complaints != "Fever" {
		t.Errorf("Expected complaints Fever, got %s", retrieved.Complaints)
	}
}
