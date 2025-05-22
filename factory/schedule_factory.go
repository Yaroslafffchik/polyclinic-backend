package factory

import (
	"errors"
	"polyclinic-backend/models"
	"strings"
)

func NewSchedule(doctorID uint, days, time, room string) (*models.Schedule, error) {
	dayList := strings.Split(days, ",")
	if len(dayList) > 3 {
		return nil, errors.New("schedule cannot exceed 3 days per week")
	}
	return &models.Schedule{
		DoctorID: doctorID,
		Days:     days,
		Time:     time,
		Room:     room,
	}, nil
}
