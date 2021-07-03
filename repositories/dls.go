package repositories

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
)

func (app *App) GetDLsByProject(projectID uint) ([]*models.DLS, error) {

	db := app.Db
	dlsArr := []*models.DLS{}
	rows, err := db.Table("dls").Select("dls.*").Joins("join dls_projects on dls_projects.dls_id = dls.id").Where("dls_projects.project_id = ?", projectID).Rows()
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
	var dls []models.DLSDetails
	query := db.Table(table).Joins("JOIN dls_projects on dls_projects.dls_id = dls.id").Joins("JOIN projects on projects.id = dls_projects.project_id")
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)
	query = query.Order("id " + dt.Order)

	if err := query.Debug().Select("dls.*, projects.name as project_name").Find(&dls).Error; err != nil {
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

func (app *App) GetDL(id string) (models.DLS, error) {
	dls := models.DLS{}
	result := app.Db.First(&dls, id)
	return dls, result.Error
}

func (app *App) GetDLProject(id string) (models.DlsProjects, error) {
	dlProject := models.DlsProjects{}
	result := app.Db.Where("dls_id = ?", id).First(&dlProject)
	return dlProject, result.Error
}

func (app *App) UpdateDL(c *gin.Context, dlData models.DLS, dlProjectId string) (uint, error) {
	dl, err := app.GetDL(c.Param("id"))
	if err != nil {
		return 0, err
	}
	app.Db.Unscoped().Where("dls_id = ?", dl.ID).Delete(&models.DlsProjects{})
	updatedProject := app.Db.Model(&dl).Updates(&dlData)
	var errMessage = updatedProject.Error
	if updatedProject.Error != nil {
		log.Print(errMessage)

	}
	projectId, err := strconv.ParseUint(dlProjectId, 10, 32)
	dlProjects := models.DlsProjects{
		DlsId:     dl.ID,
		ProjectId: uint(projectId),
	}
	createDlProjects := app.Db.Create(&dlProjects)
	if createDlProjects.Error != nil {
		log.Print(createDlProjects.Error)
	}
	return dl.ID, errMessage
}

func (app *App) CreateDL(c *gin.Context, dl models.DLS, dlProjectId string) (uint, error) {
	err := c.Bind(&dl)
	if err != nil {
		log.Print(err)
	}
	createdProject := app.Db.Create(&dl)
	var errMessage = createdProject.Error

	if createdProject.Error != nil {
		log.Print(errMessage)

	}
	projectId, err := strconv.ParseUint(dlProjectId, 10, 32)
	if err != nil {
		log.Print(errMessage)
	}
	dlProjects := models.DlsProjects{
		DlsId:     dl.ID,
		ProjectId: uint(projectId),
	}
	createDlProjects := app.Db.Create(&dlProjects)
	if createDlProjects.Error != nil {
		log.Print(errMessage)
	}
	return dl.ID, errMessage
}

func (app *App) DeleteDL(id string) (uint, error) {
	dl, err := app.GetDL(id)
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
