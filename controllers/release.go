package controllers

import (
	"fmt"
	"github.com/release-trackers/gin/cmd/jira"
	"github.com/release-trackers/gin/notifications/mails"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
)

type Milestones struct {
	BetaReleaseDate string
	DevCompletionDate string
	Project string
}

func (app *App) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "release/home", gin.H{
	})
}

func (app *App) GetListOfReleases(c *gin.Context) {
	var columnOrder string
	columnOrder = "desc"
	c.Request.ParseForm()
	order := c.PostFormMap("order")
	for _, value := range order { // Order not specified
		if value == "asc" || value == "desc" {
			columnOrder = value
		}
	}
	dtValues := models.DataTableValues{
		Offset: QueryOffset(c),
		Limit:  QueryLimit(c),
		Search: SearchScope(c),
		Order:  columnOrder,
	}
	log.Printf("length %v: ", dtValues.Order)
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	println("++++++app.Name++++++")
	println(app.Name)
	releases := releaseRepsitoryHandler.GetAllReleases(c, dtValues)
	c.JSON(http.StatusOK, releases)
}


func (app *App) CreateReleaseForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	projects, err := releaseRepsitoryHandler.GetProjects(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	c.HTML(http.StatusOK, "release/create", gin.H{
		"title":    "Create release",
		"projects": projects,
	})
}
func (app *App) ReleaseTicketsForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	releases, err := releaseRepsitoryHandler.GetLatestReleases()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	c.HTML(http.StatusOK, "release/list", gin.H{
		"title":    "Create release",
		"releases": releases,
	})
}

//Not using as of now
func (app *App)ReleaseListTickets(c *gin.Context)  {
	releaseName := c.Query("release")
	log.Printf("param : %v", releaseName)
	release := &models.Release{}
	fetchRelease := app.Db.Debug().Where("name = ?", strings.TrimSpace(releaseName)).Find(release)
	if fetchRelease.Error != nil {
		log.Print(fetchRelease.Error)
	}
	jirsList := jira.GetIssuesByLabel(releaseName)
	for _, jiraTickets := range jirsList {
		releaseTickets := &models.ReleaseTickets{Key: jiraTickets.Id, Summary: jiraTickets.Summary, Type: jiraTickets.Type,
			Project: jiraTickets.Project, Status: jiraTickets.Status, ReleaseId: release.ID}
		app.Db.Create(releaseTickets)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "List found", "data": jirsList})
	return
}

func (app *App)ReleaseListTicketsByReleaseId(c *gin.Context)  {
	releaseName := c.Query("release")
	sendMail := c.Query("sendEmail")
	log.Printf("param : %v", releaseName)
	log.Printf("send email : %v", sendMail)
	release := &models.Release{}
	fetchRelease := app.Db.Debug().Where("name = ?", strings.TrimSpace(releaseName)).Find(release)
	if fetchRelease.Error != nil {
		log.Print(fetchRelease.Error)
	}
	releaseTickets := []models.ReleaseTickets{}
	records := app.Db.Debug().Where("release_id = ?", release.ID).Find(&releaseTickets)
	ticketsrows, _ := records.Rows()
	defer ticketsrows.Close()

	ticketsrr := []*models.ReleaseTickets{}
	for ticketsrows.Next() {
		tickets := &models.ReleaseTickets{}
		err := app.Db.Debug().ScanRows(ticketsrows, &tickets)
		if err != nil {
			log.Fatalln(err)
		}
		ticketsrr = append(ticketsrr, tickets)
	}
	log.Printf("release tickets : %+v", ticketsrr)
	if sendMail ==  "true" {
		mails.SendReleaseNotes(release, ticketsrr)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "List found", "data": ticketsrr})
	return
}

func (app *App) ViewReleaseForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	releases, projects, reviewers, err := releaseRepsitoryHandler.GetReleases(c)
	jiraTickets := jira.GetIssuesByLabel(releases.Name)
	var ticketsarr []string
	for _,tickets := range jiraTickets{
		ticketId := tickets.Id
		ticketsarr = append(ticketsarr,ticketId)
	}
	var milestones []*Milestones
	for _, project := range projects {
		betaRelease := GetMilestoneDates(project.BetaReleaseDate, releases, project)
		devCompletion := GetMilestoneDates(project.DevCompletionDate, releases, project)
		mileStone := &Milestones{BetaReleaseDate: betaRelease.Format("2006-01-02"), DevCompletionDate: devCompletion.Format("2006-01-02"), Project: project.Name}

		milestones = append(milestones, mileStone)
	}
	log.Printf("after loop : %v", milestones)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	jiraBaseUrl := os.Getenv("JIRA_BASE_URL")+"browse/"
	c.HTML(http.StatusOK, "release/view", gin.H{
		"title":     "View release",
		"projects":  projects,
		"releases":  releases,
		"reviewers": reviewers,
		"tickets":	ticketsarr,
		"milestones" : milestones,
		"jiraurl" : jiraBaseUrl,
	})
}

func GetMilestoneDates(days string, release models.Release, project *models.Project) time.Time {
	daysToSubtract, err := strconv.Atoi(days)
	if err != nil {
		log.Println(err)
	}
	releaseDate := release.TargetDate.AddDate(0, 0, -daysToSubtract).Truncate(24 * time.Hour)
	return releaseDate
}

func (app *App) CreateRelease(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
		return
	}
	if c.Request.Method != http.MethodPost {
		c.Writer.Header().Set("Allow", http.MethodPost)
		http.Error(c.Writer, "Method Not Allowed", 405)
		return
	}

	title := c.Request.PostForm.Get("name")
	release_type := c.Request.PostForm.Get("type")
	target_date := c.Request.PostForm.Get("target_date")
	owner := c.Request.PostForm.Get("owner")
	projectIds := c.Request.PostForm["projects"]
	convertedProjectIds := app.covertStringToIntArray(projectIds)
	target_format, err := time.Parse(time.RFC3339, target_date+"T15:04:05Z")
	if err != nil {
		fmt.Println(err)
	}
	var release = models.Release{
		Name:       title,
		Type:       release_type,
		TargetDate: target_format,
		Owner:      owner,
	}
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	createReleaseData, err := releaseRepsitoryHandler.CreateRelease(c, release, convertedProjectIds)
	session := sessions.Default(c)
	if err != nil {
		log.Print(err)
		session.AddFlash(err, "error")
	}
	if createReleaseData != 0 {
		session.AddFlash("Release created successfully", "success")
	}
	session.Save()
	c.Redirect(http.StatusFound, "/release/index")
}

func (app *App) covertStringToIntArray(projectIds []string) []int {
	var convertedProjectIds = []int{}
	for _, i := range projectIds {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		convertedProjectIds = append(convertedProjectIds, j)
	}
	return convertedProjectIds
}

func (app *App) GetProjectReviewerList(c *gin.Context) {
	projects := c.Query("ids")
	s := strings.Split(projects, ",")
	convertedProjectIds := app.covertStringToIntArray(s)
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	revList, err := releaseRepsitoryHandler.GetReviewers(c, convertedProjectIds)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"status": "failed", "message": "No list found"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "List found", "data": revList})
	return
}

//func (app *App) GetAccessToken(c *gin.Context)  {
//	err := c.Request.ParseForm()
//	if err != nil {
//		http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
//		return
//	}
//	code := c.Request.PostForm.Get("code")
//	fmt.Printf("code entered %v", code)
//	//tokenSession := bitbucket.GetAccessToken(code)
//	tokenSession := bitbucket.GetAccessToken(c)
//	fmt.Printf("access token controller %v", tokenSession.AccessToken)
//	session := sessions.Default(c)
//	session.Set("access_token", tokenSession.AccessToken)
//	session.Save()
//	//session.CreateBranch(c,"hotfix", "bug_fix_13")
//	//c.JSON(http.StatusOK, gin.H{"status": "success", "message": "access token", "data":tokenSession.AccessToken})
//	//return
//	c.Redirect(http.StatusFound, "/release/index")
//}

//func (app *App)GetAuthCode(c *gin.Context)  {
//	code := c.Query("code")
//	c.HTML(http.StatusOK, "oauth/create", gin.H{
//		"code" : code,
//	})
//}
