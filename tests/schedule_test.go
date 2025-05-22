package tests

import (
	"polyclinic-backend/db"
	"polyclinic-backend/factory"
	"polyclinic-backend/models"
	"testing"
)

func TestCreateSchedule(t *testing.T) {
	db.Init()
	schedule, err := factory.NewSchedule(1, "Mon,Wed", "09:00-12:00", "101")
	if err != nil {
		t.Fatalf("Failed to create schedule: %v", err)
	}
	db.DB.Create(schedule)

	var retrieved models.Schedule
	if err := db.DB.Preload("Doctor").First(&retrieved, schedule.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve schedule: %v", err)
	}
	if retrieved.Days != "Mon,Wed" {
		t.Errorf("Expected days Mon,Wed, got %s", retrieved.Days)
	}
}
