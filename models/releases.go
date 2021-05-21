package models

import (
	"gorm.io/gorm"
	"time"
)

//Release struct declaration
type Release struct {
	gorm.Model
	Name          		string `gorm:"type:varchar(100);index:,unique"`
	TargetDate    		time.Time `json:"target_date"`
	Type          		string `json:"type"`
	Owner           	string  `json:"owner"`
}