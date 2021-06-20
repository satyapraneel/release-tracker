package models

import "gorm.io/gorm"

const QA = "qa"
const DEV = "dev"

//User struct declaration
type DLS struct {
	gorm.Model
	Name   string `gorm:"type:varchar(100)"`
	Email  string `gorm:"type:varchar(100)"`
	DlType string `gorm:"type:varchar(10)"`
}
