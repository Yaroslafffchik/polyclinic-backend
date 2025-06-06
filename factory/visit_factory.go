package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
)

func NewVisit(patientID, doctorID uint, date, complaints, diagnosis, prescription string, sickLeave bool, sickLeaveDuration int) (*models.Visit, error) {
	if patientID == 0 {
		return nil, errors.New("patient ID is required")
	}
	if doctorID == 0 {
		return nil, errors.New("doctor ID is required")
	}
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(date) {
		return nil, errors.New("date must be in YYYY-MM-DD format")
	}
	if sickLeave && sickLeaveDuration <= 0 {
		return nil, errors.New("sick leave duration must be greater than 0 if sick leave is issued")
	}
	return &models.Visit{
		PatientID:         patientID,
		DoctorID:          doctorID,
		Date:              date,
		Complaints:        complaints,
		Diagnosis:         diagnosis,
		Prescription:      prescription,
		SickLeave:         sickLeave,
		SickLeaveDuration: sickLeaveDuration,
	}, nil
}
