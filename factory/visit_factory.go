package factory

import (
	"errors"
	"polyclinic-backend/models"
	"regexp"
)

func NewVisit(patientID, doctorID uint, date, complaints, diagnosis, prescription string, sickLeave bool) (*models.Visit, error) {
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(date) {
		return nil, errors.New("date must be in YYYY-MM-DD format")
	}
	return &models.Visit{
		PatientID:    patientID,
		DoctorID:     doctorID,
		Date:         date,
		Complaints:   complaints,
		Diagnosis:    diagnosis,
		Prescription: prescription,
		SickLeave:    sickLeave,
	}, nil
}
