package models

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	Action    string `gorm:"type:varchar(100)"`
	Timestamp string `gorm:"type:timestamp;default:current_timestamp"`
	Duration  string `gorm:"type:varchar(50)"`
}
