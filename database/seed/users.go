package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
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
