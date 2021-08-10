package models

import "github.com/jinzhu/gorm"

type DataTableValues struct {
	Offset int
	Limit  int
	Search func(db2 *gorm.DB) *gorm.DB
	Order  string
}
