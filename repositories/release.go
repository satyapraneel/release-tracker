package repositories

import (
	"log"
	"strings"
	"time"

	"github.com/release-trackers/gin/cmd/bitbucket"
	"github.com/release-trackers/gin/cmd/jira"
	"github.com/release-trackers/gin/notifications/mails"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
)

// NewReleaseHandler ..
func NewReleaseHandler(app *cmd.Application) *App {
	return &App{app}
}

func (app *App) CreateRelease(c *gin.Context, release models.Release, projectIds []int) (uint, error) {
	//set bitbucket access token in session
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
		log.Printf("projectIds : %v", projectId)
		app.Db.Create(&models.ReleaseProject{
			ReleaseId: release.ID, ProjectId: uint(projectId),
		})
		fetchProject := app.Db.Debug().Where("id = ?", projectId).Find(project)
		if fetchProject.Error != nil {
			log.Print(errMessage)
		}
		log.Printf("brandName : %v", project.RepoName)
		log.Printf("release : %+v, %v", release.Name, release.Type)
		dlsList, _ := app.GetDLsByProject(project.ID)
		go mails.SendReleaseCreatedMail(&release, project, dlsList)
		reviewerUserNames := app.GetReviewerUserNames(c, project.ReviewerList)
		go bitbucket.CreateBranch(app.Db, release, reviewerUserNames, project.RepoName)
		project = &models.Project{}
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
		Total:    total,
		Filtered: filtered,
		Data:     release,
	}

	return result

}

func (app *App) GetReleaseProjects(release models.Release) ([]*models.Project, []string, error) {
	db := app.Db
	projects := []models.ReleaseProject{}
	projectRecords := db.Where("release_id = ?", release.ID).Find(&projects)
	projrows, err := projectRecords.Rows()
	projectArr := []*models.Project{}
	var reviewers []string
	for projrows.Next() {
		releaseProject := &models.ReleaseProject{}
		project := &models.Project{}
		err := db.ScanRows(projrows, releaseProject)
		result := app.Db.Where("status = ?", "1").First(project, releaseProject.ProjectId)
		if err != nil {
			log.Fatalln(err)
		}
		if result.Error != nil {
			continue
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

func (app *App) GetReviewerUserNames(c *gin.Context, reviewerList string) []string {
	revi := strings.Split(reviewerList, ",")
	db := app.Db
	projectRecords := db.Table("reviewers").Select("user_name").Where("email in (?)", revi)
	rows, _ := projectRecords.Rows()
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		reviewer := &models.Reviewers{}
		err := db.Debug().ScanRows(rows, &reviewer)
		if err != nil {
			log.Fatalln(err)
		}
		usernames = append(usernames, reviewer.UserName)
	}
	return usernames
}

func (app *App) GetLatestReleases() ([]models.Release, error) {
	releases := []models.Release{}
	releaseRecords := app.Db.Debug().Table("releases").Where("status = ?", 1)
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

func (app *App) CloseRelease() (int, error) {
	todaysDate := time.Now()
	today := todaysDate.Format("2006-01-02")
	updatedRecord := app.Db.Debug().Model(&models.Release{}).Where("target_date < ?", today).Update("status", 0)
	if updatedRecord.Error != nil {
		println(updatedRecord.Error)
		return 0, updatedRecord.Error
	}
	return 1, nil
}

func (app *App) UpdateJiraTicketsToDB(jirsList []*jira.JiraTickets, releaseId uint) {
	for _, jiraTickets := range jirsList {
		releaseTickets := &models.ReleaseTickets{Key: jiraTickets.Id, Summary: jiraTickets.Summary, Type: jiraTickets.Type,
			Project: jiraTickets.Project, Status: jiraTickets.Status, ReleaseId: releaseId}
		createdReleaseTickets := app.Db.Create(releaseTickets)
		var errMessage = createdReleaseTickets.Error
		if createdReleaseTickets.Error != nil {
			log.Print(errMessage)
		}
	}
}
