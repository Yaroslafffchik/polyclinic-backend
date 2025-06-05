package models

import "gorm.io/gorm"

type Section struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Address string `gorm:"type:varchar(200);not null" json:"address"`
}
