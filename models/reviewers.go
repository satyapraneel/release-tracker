package models

import "github.com/jinzhu/gorm"

//BitbucketOAuth struct declaration
type Reviewers struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);index:,unique"`
	Name     string `json:"name"`
	UserName string `gorm:"type:varchar(100);index:,unique" json:"username"`
}

type ReviewerResult struct {
	Total    int64       `json:"recordsTotal"`
	Filtered int64       `json:"recordsFiltered"`
	Data     []Reviewers `json:"data"`
}
