package models

import (
	"time"

	"gorm.io/gorm"
)

//Release struct declaration
type Release struct {
	gorm.Model
	Name       string    `gorm:"type:varchar(100);index:,unique"`
	TargetDate time.Time `json:"target_date"`
	Type       string    `json:"type"`
	Owner      string    `json:"owner"`
	RestrictionId uint	`json:"restriction_id"`
}

type DataResult struct {
	Total    int64     `json:"recordsTotal"`
	Filtered int64     `json:"recordsFiltered"`
	Data     []Release `json:"data"`
}
