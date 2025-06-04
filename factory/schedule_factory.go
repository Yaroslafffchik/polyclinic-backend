package factory

import (
	"fmt"
	"polyclinic-backend/models"
)

func NewSchedule(doctorID uint, days, time, room string) (*models.Schedule, error) {
	if doctorID == 0 {
		return nil, fmt.Errorf("doctor ID is required")
	}
	if days == "" {
		return nil, fmt.Errorf("days are required")
	}
	if time == "" {
		return nil, fmt.Errorf("time is required")
	}
	if room == "" {
		return nil, fmt.Errorf("room is required")
	}

	return &models.Schedule{
		DoctorID: doctorID,
		Days:     days,
		Time:     time,
		Room:     room,
	}, nil
}
