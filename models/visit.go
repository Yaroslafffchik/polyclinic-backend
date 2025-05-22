package models

import "gorm.io/gorm"

type Visit struct {
	gorm.Model
	PatientID    uint
	Patient      Patient
	DoctorID     uint
	Doctor       Doctor
	Date         string `gorm:"type:varchar(10)"`
	Complaints   string `gorm:"type:text"`
	Diagnosis    string `gorm:"type:text"`
	Prescription string `gorm:"type:text"`
	SickLeave    bool
}
