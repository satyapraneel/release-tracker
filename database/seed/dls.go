package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
)

func CreateDls(db *gorm.DB,
	name string,
	email string,
) error {
	return db.Create(&models.
		DLS{
		Name:  name,
		Email: email,
	}).Error
}
