package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	FullName       string `gorm:"type:varchar(100);not null"`
	Category       string `gorm:"type:varchar(50)"`
	Experience     int
	BirthDate      string `gorm:"type:varchar(10)"`
	Specialization string `gorm:"type:varchar(50)"`
	SectionID      uint
	Section        Section
}
