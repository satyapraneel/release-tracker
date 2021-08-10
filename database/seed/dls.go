package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/release-trackers/gin/models"
)

func CreateDls(db *gorm.DB,
	email string,
	dlType string,
) error {
	return db.Create(&models.
		DLS{
		Email:  email,
		DlType: dlType,
	}).Error
}
