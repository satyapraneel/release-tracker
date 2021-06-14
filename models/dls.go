package models

import "gorm.io/gorm"

//User struct declaration
type DLS struct {
	gorm.Model
	Name  string `gorm:"type:varchar(100)"`
	Email string `gorm:"type:varchar(100);unique_index"`
}
