package repositories

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
	"go.elastic.co/apm/module/apmgorm"
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

func (app *App) GetAllProjects(c *gin.Context, dt models.DataTableValues) models.ProjectResult {
	table := "projects"
	// db := app.Db
	db := apmgorm.WithContext(c.Request.Context(), app.Db)
	var total, filtered int64
	var project []models.Project
	query := db.Table(table)
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)
	query = query.Order("id " + dt.Order)

	if err := query.Find(&project).Error; err != nil {
		return models.ProjectResult{
			Total:    0,
			Filtered: 0,
			Data:     project,
		}
	}

	// Filtered data count
	query.Table(table).Count(&filtered)

	// Total data count
	db.Table(table).Count(&total)

	return models.ProjectResult{
		Total:    total,
		Filtered: filtered,
		Data:     project,
	}

}

func (app *App) CreateProject(c *gin.Context, project models.Project) (uint, error) {
	err := c.Bind(&project)
	if err != nil {
		log.Print(err)
	}
	createdProject := app.Db.Create(&project)
	var errMessage = createdProject.Error

	if createdProject.Error != nil {
		log.Print(errMessage)

	}
	return project.ID, errMessage
}

func (app *App) GetProject(c *gin.Context) (models.Project, error) {
	id := c.Param("id")
	project := models.Project{}
	result := app.Db.First(&project, id)
	return project, result.Error
}

func (app *App) UpdateProject(c *gin.Context, projectData models.Project) (uint, error) {
	project, err := app.GetProject(c)
	if err != nil {
		return 0, err
	}
	updatedProject := app.Db.Debug().Model(&project).Updates(&projectData)
	var errMessage = updatedProject.Error
	if updatedProject.Error != nil {
		log.Print(errMessage)
	}
	//since empty string was not updating
	if projectData.ReviewerList == "" && updatedProject.Error == nil {
		updatedProject = app.Db.Debug().Model(&project).Updates(map[string]interface{}{"reviewer_list": ""})
	}
	errMessage = updatedProject.Error
	if updatedProject.Error != nil {
		log.Print(errMessage)
	}
	return project.ID, errMessage
}

func (app *App) GetAllProjectsList(c *gin.Context) ([]*models.Project, error) {
	db := app.Db
	var projectList []models.Project
	records := db.Find(&projectList)
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
		err := db.ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr = append(projectArr, project)
	}
	return projectArr, err
}
