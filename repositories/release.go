package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/models"
	"log"
)

var (
	db = database.InitConnection()
)

type App struct {
	*cmd.Application
}

// NewReleaseHandler ..
func NewReleaseHandler(app *cmd.Application) *App {
	return &App{app}
}

func (a *App) CreateRelease(c *gin.Context, release models.Release, projectIds []int) (uint, error) {
	err := c.Bind(&release)
	if err != nil {
		log.Print(err)
	}
	createdRelease := db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		log.Print(errMessage)

	}

	for _, projectId := range projectIds {
		db.Create(&models.ReleaseProject{
			ReleaseId: release.ID, ProjectId: uint(projectId),
		})
	}
	return release.ID, errMessage
}

func (app *App) GetAllReleases(c *gin.Context, dt models.DataTableValues)  (models.DataResult) {
	table := "releases"
	var total, filtered int64
	var release []models.Release
	query := db.Table(table)
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)

	if err := query.Find(&release).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
	}

	// Filtered data count
	query.Table(table).Count(&filtered)

	// Total data count
	db.Table(table).Count(&total)

	result := models.DataResult{
		total,
		filtered,
		release,
	}

	return result

}

func (app *App) getReleaseProjects(release *models.Release, err error)  ([]*models.Project) {
	projects := []models.ReleaseProject{}
	log.Printf("release Id : %+v", release.ID)
	projectRecords := db.Debug().Where("release_id = ?", release.ID).Find(&projects)
	projrows, err := projectRecords.Rows()
	log.Printf("project %+v\n", projrows)
	projectArr := []*models.Project{}

	for projrows.Next() {
		project:= &models.Project{}
		err := db.Debug().ScanRows(projrows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr = append(projectArr, project)
	}

	return projectArr
}

func (app *App) GetProjects(c *gin.Context) ([]*models.Project, error) {
	var project []models.Project
	records := db.Debug().Find(&project)
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
		project:= &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr=append(projectArr, project)
	}
	return projectArr, err
}


func GetReviewers(c *gin.Context, projectIds []int) ([]string , error) {
	projectRecords := db.Table("projects").Select("reviewer_list").Where("id in (?)", projectIds)
	rows, err := projectRecords.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		project:= &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		reviewers=append(reviewers, project.ReviewerList)
	}
	return reviewers, err
}
