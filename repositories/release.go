package repositories

import (
	"github.com/release-trackers/gin/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
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

func (app *App) GetAllReleases(c *gin.Context) ([]*models.Release, error) {
	var release []models.Release
	records := db.Debug().Find(&release)
	if records.Error != nil {
		log.Fatalln(records.Error)
	}
	log.Printf("%d rows found.", records.RowsAffected)
	rows, err := records.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	releaseArr := []*models.Release{}
	//projectArr := []*models.Project{}

	for rows.Next() {
		release := &models.Release{}
		err := db.Debug().ScanRows(rows, &release)
		if err != nil {
			log.Fatalln(err)
		}
		//projectArr = getReleaseProjects(release, err)
		releaseArr = append(releaseArr, release)
	}
	return releaseArr, err
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
		//log.Printf("%+v\n", project)
		projectArr=append(projectArr, project)
	}
	return projectArr, err
}


func GetReviewers(c *gin.Context, projectIds []int) ([]string , error) {
	log.Printf("release Id : %+v", projectIds)
	//projectRecords := db.Debug().Exec("select reviewer_list from projects where project_id IN (?)", projectIds)
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
		//log.Printf("%+v\n", project)
		reviewers=append(reviewers, project.ReviewerList)
	}
	return reviewers, err
}
