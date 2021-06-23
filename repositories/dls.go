package repositories

import (
	"log"

	"github.com/gin-gonic/gin"
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

func (app *App) GetAllDls(dt models.DataTableValues) models.DLSResult {
	table := "dls"
	db := app.Db
	var total, filtered int64
	var dls []models.DLS
	query := db.Table(table)
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)
	query = query.Order("id " + dt.Order)

	if err := query.Find(&dls).Error; err != nil {
		return models.DLSResult{
			Total:    0,
			Filtered: 0,
			Data:     dls,
		}
	}

	// Filtered data count
	query.Table(table).Count(&filtered)

	// Total data count
	db.Table(table).Count(&total)

	return models.DLSResult{
		Total:    total,
		Filtered: filtered,
		Data:     dls,
	}
}

func (app *App) GetDL(c *gin.Context) (models.DLS, error) {
	id := c.Param("id")
	dls := models.DLS{}
	result := app.Db.First(&dls, id)
	return dls, result.Error
}

func (app *App) UpdateDL(c *gin.Context, dlData models.DLS) (uint, error) {
	dl, err := app.GetDL(c)
	if err != nil {
		return 0, err
	}
	updatedProject := app.Db.Model(&dl).Updates(&dlData)
	var errMessage = updatedProject.Error
	if updatedProject.Error != nil {
		log.Print(errMessage)

	}
	return dl.ID, errMessage
}

func (app *App) CreateDL(c *gin.Context, dl models.DLS) (uint, error) {
	err := c.Bind(&dl)
	if err != nil {
		log.Print(err)
	}
	createdProject := app.Db.Create(&dl)
	var errMessage = createdProject.Error

	if createdProject.Error != nil {
		log.Print(errMessage)

	}
	return dl.ID, errMessage
}

func (app *App) DeleteDL(c *gin.Context) (uint, error) {
	dl, err := app.GetDL(c)
	if err != nil {
		return 0, err
	}
	app.Db.Unscoped().Where("dls_id = ?", dl.ID).Delete(&models.DlsProjects{})
	// app.Db.Unscoped().Delete(&models.DlsProjects{}, dl.ID)
	app.Db.Unscoped().Delete(&dl)
	return 1, nil
}

func (app *App) GetAllDLsList(c *gin.Context) ([]*models.DLS, error) {
	db := app.Db
	var dlsList []models.DLS
	records := db.Find(&dlsList)
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
		err := db.ScanRows(rows, &dls)
		if err != nil {
			log.Fatalln(err)
		}
		dlsArr = append(dlsArr, dls)
	}
	return dlsArr, err
}
