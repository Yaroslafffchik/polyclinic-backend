package models

import "gorm.io/gorm"

type Section struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Address string `gorm:"type:varchar(200);not null"` // Поле Address теперь обязательно
}
