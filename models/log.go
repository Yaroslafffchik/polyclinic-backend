package models

import (
	"gorm.io/gorm"
	"time"
)

type Log struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Action    string    `gorm:"type:varchar(100)"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Duration  string    `gorm:"type:varchar(50)"`
}
