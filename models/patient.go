package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	FullName        string `gorm:"type:varchar(100);not null" json:"full_name"`
	Address         string `gorm:"type:varchar(200)" json:"address"`
	Gender          string `gorm:"type:char(1)" json:"gender"`
	Age             int    `json:"age"`
	InsuranceNumber string `gorm:"type:varchar(16);unique" json:"insurance_number"`
	UserID          uint   `json:"user_id"`
	User            User   `gorm:"foreignKey:UserID"`
	DoctorID        uint   `json:"doctor_id"`
	Doctor          Doctor `gorm:"foreignKey:DoctorID"`
}
