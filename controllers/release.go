package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/models"
	repositories "github.com/release-trackers/gin/repositories"
	"log"
	"net/http"
	"strconv"
	"time"
)
type App struct {
	*cmd.Application
}

// NewReleaseHandler ..
func NewReleaseHandler(server *cmd.Application) App {
	return App{server}
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

type BlogPost struct {
	Date time.Time
}

func (post *BlogPost) formattedDate() string {
	return post.Date.Format(time.RFC822)
}

func (app App) GetListOfReleases(c *gin.Context) {
	releases, err := repositories.GetAllReleases(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	//router := gin.Default()
	//router.Delims("{[{", "}]}")
	//router.SetFuncMap(template.FuncMap{
	//	"formatAsDate": formatAsDate,
	//})
	//router.LoadHTMLFiles("ui/html/release/home.tmpl")
	c.HTML(http.StatusOK, "release/home", gin.H{
		"releases": releases,

	})
}

func (app App) CreateReleaseForm(c *gin.Context) {
	projects, err := repositories.GetProjects(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	c.HTML(http.StatusOK, "release/create", gin.H{
		"title": "Create release",
		"projects" : projects,
	})
}

func (app App) CreateRelease(c *gin.Context) {
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
	log.Printf("projectIds %v", projectIds)
	var convertedProjectIds = []int{}
	log.Printf("convertedProjectIds %v", convertedProjectIds)
	for _, i := range projectIds {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		convertedProjectIds = append(convertedProjectIds, j)
	}

	log.Printf("after loop convertedProjectIds %v", convertedProjectIds)
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
		Owner: owner,
	}

	log.Printf("%v release info", release)
	createReleaseData, err := repositories.CreateRelease(c, release, convertedProjectIds)
	if err != nil {
		log.Print(err)
	}
	log.Print(createReleaseData)
	c.Redirect(http.StatusFound, "/release/list")
}
