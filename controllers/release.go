package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
	"strconv"
)

type App struct {
	*cmd.Application
}

// NewReleaseHandler ..
func NewReleaseHandler(app *cmd.Application) *App {
	return &App{ app}
}

func (app *App) GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "release/home", gin.H{
	})
}
func (app *App) GetListOfReleases(c *gin.Context) {
	c.Request.ParseForm();
	dtValues := models.DataTableValues{
		Offset: QueryOffset(c),
		Limit: QueryLimit(c),
		Search: SearchScope(c),
	}
	log.Printf("length %v: ", c.Request.PostForm.Get("length"))
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
		"title": "Create release",
		"projects" : projects,
	})
}

func (app *App) ViewReleaseForm(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	releases, projects, reviewers, err := releaseRepsitoryHandler.GetReleases(c)
	for _, project := range projects {
		log.Printf("project in loop : %v", project.Name)
	}
	log.Printf("after loop : %v", releases.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	c.HTML(http.StatusOK, "release/view", gin.H{
		"title": "View release",
		"projects" : projects,
		"releases" : releases,
		"reviewers":reviewers,
	})
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
	if err != nil {
		log.Print(err)
	}
	log.Print(createReleaseData)
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

func (app *App) GetProjectReviewerList(c *gin.Context)  {
	projects := c.Query("ids")
	s := strings.Split(projects, ",")
	convertedProjectIds := app.covertStringToIntArray(s)
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	revList, err := releaseRepsitoryHandler.GetReviewers(c, convertedProjectIds)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"status": "failed", "message": "No list found"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "List found", "data":revList})
	return
}
