package models

import "gorm.io/gorm"

//BitbucketOAuth struct declaration
type BitbucketOAuth struct {
	gorm.Model
	Code 					string `gorm:"type:varchar(100);index:,unique"`
	AccessToken       		string `json:"access_token"`
	RefreshToken    		string `json:"refresh_token"`
	ExpiresAt    			string `json:"expires_at"`
}
