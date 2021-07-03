package models

import (
	"time"

	"gorm.io/gorm"
)

//Release struct declaration
type Release struct {
	gorm.Model
	Name               string    `gorm:"type:varchar(100);index:,unique"`
	TargetDate         time.Time `json:"target_date"`
	Type               string    `json:"type"`
	Owner              string    `json:"owner"`
	RestrictionPushId  uint      `json:"restriction_push_id"`
	RestrictionMergeId uint      `json:"restriction_merge_id"`
	Status             byte      `json:"status"`
}

type DataResult struct {
	Total    int64     `json:"recordsTotal"`
	Filtered int64     `json:"recordsFiltered"`
	Data     []Release `json:"data"`
}
