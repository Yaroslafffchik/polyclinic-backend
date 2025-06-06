package models

import "gorm.io/gorm"

type DoctorSections struct {
	gorm.Model
	DoctorID  uint    `gorm:"not null" json:"doctor_id"`
	SectionID uint    `gorm:"not null" json:"section_id"`
	Doctor    Doctor  `gorm:"foreignKey:DoctorID"`
	Section   Section `gorm:"foreignKey:SectionID"`
}
