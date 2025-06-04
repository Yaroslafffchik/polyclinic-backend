package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	FullName       string `gorm:"type:varchar(100);not null" json:"full_name"`
	Category       string `gorm:"type:varchar(50)" json:"category"`
	Experience     int    `json:"experience"`
	BirthDate      string `gorm:"type:varchar(10)" json:"birth_date"`
	Specialization string `gorm:"type:varchar(50)" json:"specialization"`
	SectionID      uint   `json:"section_id"`
	Section        Section
	UserID         uint `json:"user_id"`
	User           User `gorm:"foreignKey:UserID"`
}
