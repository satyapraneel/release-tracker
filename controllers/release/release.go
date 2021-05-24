package release

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/repositories"
	"net/http"
	"time"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func GetListOfReleases (c *gin.Context)  {
	releases, err := repositories.GetAllReleases(c);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
	}
	//router := gin.Default()
	//router.SetFuncMap(template.FuncMap{
	//	"humanDate": humanDate,
	//})
	//router.LoadHTMLFiles("./ui/html/release/home.tmpl")
	c.HTML(http.StatusOK, "release/home.tmpl", gin.H{
		"releases": releases,
		"formatAsDate": formatAsDate,
	})
}



func CreateRelease(c *gin.Context) {

}
