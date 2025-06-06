package models

import "gorm.io/gorm"

type Visit struct {
	gorm.Model
	PatientID         uint `json:"patient_id"`
	Patient           Patient
	DoctorID          uint `json:"doctor_id"`
	Doctor            Doctor
	Date              string `gorm:"type:varchar(10)" json:"date"`
	Complaints        string `gorm:"type:text" json:"complaints"`
	Diagnosis         string `gorm:"type:text" json:"diagnosis"`
	Prescription      string `gorm:"type:text" json:"prescription"`
	SickLeave         bool   `json:"sick_leave"`
	SickLeaveDuration int    `json:"sick_leave_duration"`
}
