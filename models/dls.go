package models

import "gorm.io/gorm"

const QA = "qa"
const DEV = "dev"

//User struct declaration
type DLS struct {
	gorm.Model
	Email  string `gorm:"type:varchar(100)"`
	DlType string `gorm:"type:varchar(10)"`
}

type DLSResult struct {
	Total    int64 `json:"recordsTotal"`
	Filtered int64 `json:"recordsFiltered"`
	Data     []DLS `json:"data"`
}
