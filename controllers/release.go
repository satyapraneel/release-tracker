package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
)

type App struct {
	*cmd.Application
}

// NewReleaseHandler ..
func NewReleaseHandler(app *cmd.Application) *App {
	return &App{app}
}

func (app *App) GetListOfReleases(c *gin.Context) {
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	println("++++++app.Name++++++")
	println(app.Name)
	releases, err := releaseRepsitoryHandler.GetAllReleases(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	//router := gin.Default()
	//router.SetFuncMap(template.FuncMap{
	//	"releases": releases,
	//})
	//router.LoadHTMLFiles("./ui/html/release/home.tmpl")
	c.HTML(http.StatusOK, "release/home", gin.H{
		"releases": releases,
	})
}

func (app *App) CreateReleaseForm(c *gin.Context) {
	c.HTML(http.StatusOK, "release/create", gin.H{
		"title": "Create release",
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
	//layout := "2006-01-02 15:04:05"
	target_format, err := time.Parse(time.RFC3339, target_date+"T15:04:05Z")
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("form date string ", c.Request.PostForm.Get("target_date"))
	log.Printf("form title ", target_format)
	var release = models.Release{
		Name:       title,
		Type:       release_type,
		TargetDate: target_format,
		Owner:      owner,
	}
	releaseRepsitoryHandler := repositories.NewReleaseHandler(app.Application)
	log.Printf("%v release info", release)
	createReleaseData, err := releaseRepsitoryHandler.CreateRelease(c, release)
	if err != nil {
		log.Print(err)
	}
	log.Print(createReleaseData)
	c.Redirect(http.StatusFound, "/release/list")
}
