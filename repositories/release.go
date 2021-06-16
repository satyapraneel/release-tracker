package repositories

import (
	"log"

	"github.com/release-trackers/gin/cmd/bitbucket"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/notifications/mails"
)

// NewReleaseHandler ..
func NewReleaseHandler(app *cmd.Application) *App {
	return &App{app}
}

func (app *App) CreateRelease(c *gin.Context, release models.Release, projectIds []int) (uint, error) {
	err := c.Bind(&release)
	if err != nil {
		log.Print(err)
	}
	createdRelease := app.Db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		log.Print(errMessage)

	}
	project := &models.Project{}
	for _, projectId := range projectIds {
		app.Db.Create(&models.ReleaseProject{
			ReleaseId: release.ID, ProjectId: uint(projectId),
		})
		fetchProject := app.Db.Debug().Where("id = ?", projectId).Find(project)
		if fetchProject.Error != nil {
			log.Print(errMessage)
		}
		log.Printf("reviewws : %v", project.ReviewerList)
		// cmd.TriggerMail(project.ReviewerList, release.Name, project.Name)
		go mails.SendReleaseCreatedMail(&release, project)
		// mails.SendReleaseCreatedMail(&release, project)

		bitbucket.CreateBranch(c, release.Type, release.Name, project.ReviewerList)
	}

	return release.ID, errMessage
}

func (app *App) GetAllReleases(c *gin.Context, dt models.DataTableValues) models.DataResult {
	table := "releases"
	db := app.Db
	var total, filtered int64
	var release []models.Release
	query := db.Table(table)
	query = query.Offset(dt.Offset)
	query = query.Limit(dt.Limit)
	query = query.Scopes(dt.Search)
	query = query.Order("id " + dt.Order)

	if err := query.Debug().Find(&release).Error; err != nil {
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

func (app *App) GetReleaseProjects(release models.Release) ([]*models.Project, []string, error) {
	db := app.Db
	projects := []models.ReleaseProject{}
	log.Printf("release Id : %+v", release.ID)
	projectRecords := db.Debug().Where("release_id = ?", release.ID).Find(&projects)
	projrows, err := projectRecords.Rows()
	projectArr := []*models.Project{}
	var reviewers []string
	for projrows.Next() {
		releaseProject := &models.ReleaseProject{}
		project := &models.Project{}
		err := db.Debug().ScanRows(projrows, releaseProject)
		app.Db.Debug().First(project, releaseProject.ProjectId)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v\n", project.Name)
		projectArr = append(projectArr, project)
		reviewers = append(reviewers, project.ReviewerList)
	}

	return projectArr, reviewers, err
}

func (app *App) GetProjects(c *gin.Context) ([]*models.Project, error) {
	db := app.Db
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
		project := &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		projectArr = append(projectArr, project)
	}
	return projectArr, err
}

func (app *App) GetReviewers(c *gin.Context, projectIds []int) ([]string, error) {
	db := app.Db
	projectRecords := db.Table("projects").Select("reviewer_list").Where("id in (?)", projectIds)
	rows, err := projectRecords.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		project := &models.Project{}
		err := db.Debug().ScanRows(rows, &project)
		if err != nil {
			log.Fatalln(err)
		}
		reviewers = append(reviewers, project.ReviewerList)
	}
	return reviewers, err
}

func (app *App) GetReleases(c *gin.Context) (models.Release, []*models.Project, []string, error) {
	id := c.Param("id")
	release := models.Release{}
	app.Db.First(&release, id)
	log.Printf("id : %v", release.Name)
	releaseProjects, reviewerList, errs := app.GetReleaseProjects(release)
	return release, releaseProjects, reviewerList, errs
}

func (app *App) GetLatestReleases() ([]models.Release, error) {
	releases := []models.Release{}
	releaseRecords := app.Db.Table("releases").Where("id IN (?)", app.Db.Table("releases").Select("MAX(id)").Group("type"))
	releaseRows, err := releaseRecords.Rows()
	if err != nil {
		log.Fatalln(err)
	}
	defer releaseRows.Close()

	for releaseRows.Next() {
		release := models.Release{}
		err := app.Db.Debug().ScanRows(releaseRows, &release)
		if err != nil {
			log.Fatalln(err)
		}
		releases = append(releases, release)
	}
	return releases, err
}
