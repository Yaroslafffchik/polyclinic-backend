package factory

import (
	"errors"
	"polyclinic-backend/models"
	"strings"
)

func NewSchedule(doctorID, sectionID uint, days, time, room string) (*models.Schedule, error) {
	if doctorID == 0 {
		return nil, errors.New("doctor ID must be provided")
	}
	if sectionID == 0 {
		return nil, errors.New("section ID must be provided")
	}
	if time == "" {
		return nil, errors.New("time cannot be empty")
	}
	if room == "" {
		return nil, errors.New("room cannot be empty")
	}

	// Валидация дней недели
	validDays := map[string]bool{
		"Пн": true, "Вт": true, "Ср": true, "Чт": true, "Пт": true, "Сб": true, "Вс": true,
	}
	dayList := strings.Split(strings.ReplaceAll(days, " ", ""), ",")
	if len(dayList) > 3 {
		return nil, errors.New("cannot schedule more than 3 days per week")
	}
	for _, day := range dayList {
		if !validDays[day] {
			return nil, errors.New("invalid day of the week: " + day)
		}
	}

	return &models.Schedule{
		DoctorID:  doctorID,
		SectionID: sectionID,
		Days:      days,
		Time:      time,
		Room:      room,
	}, nil
}
