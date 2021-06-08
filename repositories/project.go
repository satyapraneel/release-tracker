package repositories

import (
	"log"

	"github.com/release-trackers/gin/models"
)

func (app *App) GetProjectsByIds(projectId []int) ([]*models.Project, error) {
	db := app.Db
	var project []models.Project
	records := db.Debug().Find(&project, projectId)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	//log.Printf("%d project rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	projectArr := []*models.Project{}
	for rows.Next() {
		project := &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr = append(projectArr, project)
	}
	return projectArr, err
}
