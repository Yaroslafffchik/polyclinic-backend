package models

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model
	DoctorID  uint    `gorm:"not null" json:"doctor_id"`
	SectionID uint    `gorm:"not null" json:"section_id"`
	Days      string  `gorm:"type:varchar(50);not null" json:"days"`
	Time      string  `gorm:"type:varchar(50);not null" json:"time"`
	Room      string  `gorm:"type:varchar(50);not null" json:"room"`
	Doctor    Doctor  `gorm:"foreignKey:DoctorID"`
	Section   Section `gorm:"foreignKey:SectionID"`
}
