package models

import (
	"gorm.io/gorm"
	"time"
)

type Patient struct {
	gorm.Model
	LastName        string    `gorm:"type:varchar(50);not null" json:"last_name"`
	FirstName       string    `gorm:"type:varchar(50);not null" json:"first_name"`
	MiddleName      string    `gorm:"type:varchar(50)" json:"middle_name"`
	Address         string    `gorm:"type:varchar(200)" json:"address"`
	Gender          string    `gorm:"type:char(1)" json:"gender"`
	Age             int       `json:"age"`
	InsuranceNumber string    `gorm:"type:varchar(16);unique" json:"insurance_number"`
	CardCreatedAt   time.Time `gorm:"default:current_timestamp" json:"card_created_at"`
	UserID          uint      `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID"`
	DoctorID        uint      `json:"doctor_id"`
	Doctor          Doctor    `gorm:"foreignKey:DoctorID"`
}
