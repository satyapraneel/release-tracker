package models

import (
	"gorm.io/gorm"
	"time"
)

//Release struct declaration
type Release struct {
	gorm.Model
	Name          	string `json:"name"`
	TargetDate    	time.Time `json:"target_date"`
	Type          	string `json:"type"`
	Project			string 	`json:"project"`
	Owner           string  `json:"owner"`
}