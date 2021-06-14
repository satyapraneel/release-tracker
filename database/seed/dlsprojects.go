package seed

import (
	"github.com/release-trackers/gin/models"
	"gorm.io/gorm"
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
