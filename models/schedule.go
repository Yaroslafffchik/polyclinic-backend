package models

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model
	DoctorID uint
	Doctor   Doctor
	Days     string `gorm:"type:varchar(100)"`
	Time     string `gorm:"type:varchar(50)"`
	Room     string `gorm:"type:varchar(10)"`
}
