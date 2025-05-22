package models

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	FullName        string `gorm:"type:varchar(100);not null"`
	Address         string `gorm:"type:varchar(200)"`
	Gender          string `gorm:"type:char(1)"`
	Age             int
	InsuranceNumber string `gorm:"type:varchar(16);unique"`
}
