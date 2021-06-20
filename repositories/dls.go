package repositories

import (
	"log"

	"github.com/release-trackers/gin/models"
)

func (app *App) GetDLsByProject(projectID uint, dlType string) ([]*models.DLS, error) {

	db := app.Db
	dlsArr := []*models.DLS{}
	rows, err := db.Table("dls").Select("dls.*").Joins("join dls_projects on dls_projects.dls_id = dls.id").Where("dls_projects.project_id = ?", projectID).Where("dl_type = ?", dlType).Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		dls := &models.DLS{}
		err := db.Debug().ScanRows(rows, &dls)
		if err != nil {
			log.Fatalln(err)
		}
		dlsArr = append(dlsArr, dls)
	}
	return dlsArr, err
}
