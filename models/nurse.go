package models

import "gorm.io/gorm"

type Nurse struct {
	gorm.Model
	LastName   string  `gorm:"type:varchar(50);not null" json:"last_name"`
	FirstName  string  `gorm:"type:varchar(50);not null" json:"first_name"`
	MiddleName string  `gorm:"type:varchar(50)" json:"middle_name"`
	SectionID  uint    `json:"section_id"`
	Section    Section `gorm:"foreignKey:SectionID"`
	UserID     uint    `json:"user_id"`
	User       User    `gorm:"foreignKey:UserID"`
}
