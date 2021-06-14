package models

import "gorm.io/gorm"

//BitbucketOAuth struct declaration
type Reviewers struct {
	gorm.Model
	Email 				string `gorm:"type:varchar(100);index:,unique"`
	Name       			string `json:"name"`
	UserName    		string `gorm:"type:varchar(100);index:,unique" json:"username"`
}
