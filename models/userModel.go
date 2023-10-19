package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	//поля и типы
	Email    string `gorm:"unique"`
	Password string
}
