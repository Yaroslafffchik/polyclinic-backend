package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique"`
	Password string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:varchar(20)"`
}
