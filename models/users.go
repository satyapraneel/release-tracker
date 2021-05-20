package models

import "gorm.io/gorm"

//User struct declaration
type Users struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string `json:"Password"`
}
