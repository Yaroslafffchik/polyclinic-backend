package models

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	LastName       string    `gorm:"type:varchar(50);not null" json:"last_name"`
	FirstName      string    `gorm:"type:varchar(50);not null" json:"first_name"`
	MiddleName     string    `gorm:"type:varchar(50)" json:"middle_name"`
	Category       string    `gorm:"type:varchar(50)" json:"category"`
	Experience     int       `json:"experience"`
	BirthDate      string    `gorm:"type:varchar(10)" json:"birth_date"`
	Specialization string    `gorm:"type:varchar(50)" json:"specialization"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID"`
	Sections       []Section `gorm:"many2many:doctor_sections"`
}
