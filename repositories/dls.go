package repositories

import (
	"log"

	"github.com/release-trackers/gin/models"
)

func (app *App) getAllDLs() ([]*models.DLS, error) {
	db := app.Db
	var dls []models.DLS
	records := db.Debug().Find(&dls)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	//log.Printf("%d project rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	dlsArr := []*models.DLS{}
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
