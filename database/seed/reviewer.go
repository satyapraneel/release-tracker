package seed

import (
"github.com/release-trackers/gin/models"
"gorm.io/gorm"
)

func CreateReviwers(db *gorm.DB, name string,
	email string, user_name string,
) error {
	return db.Create(&models.
	Reviewers{
		Name:          name,
		Email:         email,
		UserName:      user_name,

	}).Error
}

