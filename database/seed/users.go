package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/release-trackers/gin/models"
)

func CreateUser(db *gorm.DB,
	name string,
	email string,
	password string,
) error {
	return db.Create(&models.
		Users{
		Name:     name,
		Email:    email,
		Password: password,
	}).Error
}
