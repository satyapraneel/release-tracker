package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
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
