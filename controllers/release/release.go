package release

import (
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/models"
	"github.com/release-trackers/gin/repositories"
	"log"
	"net/http"
	"time"
)

func GetListOfReleases(c *gin.Context) {
	releases, err := repositories.GetAllReleases(c)
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

func CreateReleaseForm(c *gin.Context) {
	c.HTML(http.StatusOK, "release/create", gin.H{
		"title": "Create release",
	})
}

func CreateRelease(c *gin.Context) {
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
	layout := "2021-05-24 18:43:31"
	target_format, _ := time.Parse(layout, target_date)
	log.Printf("form date string ", c.Request.PostForm.Get("target_date"))
	log.Printf("form title ", target_format)
	var release = models.Release{
		Name:       title,
		Type:       release_type,
		TargetDate: target_format,
		Owner: owner,
	}

	log.Printf("%v release info", release)
	createReleaseData, err := repositories.CreateRelease(c, release)
	if err != nil {
		log.Print(err)
	}
	log.Print(createReleaseData)
	c.Redirect(http.StatusFound, "/release/list")
}
