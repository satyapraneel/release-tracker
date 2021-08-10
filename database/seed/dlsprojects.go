package seed

import (
	"github.com/jinzhu/gorm"
	"github.com/release-trackers/gin/models"
)

func CreateDlsProject(db *gorm.DB,
	projectId uint,
	dlsId uint,
) error {
	return db.Create(&models.
		DlsProjects{
		ProjectId: projectId,
		DlsId:     dlsId,
	}).Error
}
