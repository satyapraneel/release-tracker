package release

import (
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/repositories"
	"log"
	"net/http"
)

func GetListOfReleases(c *gin.Context) {
	releases, err := repositories.GetAllReleases(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		//	return
	}
	log.Printf("%v in home", releases)
	//router := gin.Default()
	//router.SetFuncMap(template.FuncMap{
	//	"releases": releases,
	//})
	//router.LoadHTMLFiles("./ui/html/release/home.tmpl")
	c.HTML(http.StatusOK, "release/home", gin.H{
		"releases": releases,
	})
}

func CreateRelease(c *gin.Context) {

}
