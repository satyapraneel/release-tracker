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

type DataResult struct {
	Total    int64    `json:"recordsTotal"`
	Filtered int64    `json:"recordsFiltered"`
	Data     []Release `json:"data"`
}

type DataTableValues struct {
	Offset  int
	Limit	int
	Search  func(db2 *gorm.DB)*gorm.DB
}