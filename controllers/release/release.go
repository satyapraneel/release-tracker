package release

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/database"
	"github.com/release-trackers/gin/models"
	"log"
	"net/http"
	"time"
)

var (
	errInvalidBody     = errors.New("Invalid request body")
	errNotExist			= errors.New("No records found")
	db = database.InitConnection()
)
func CreateUser(c *gin.Context) {
	log.Print("in create user method")
	// Get DB from Mysql Config
	release := models.Release{}
	err := c.Bind(&release)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}
	log.Print("in create user method 2")
	release.ID=3
	release.Name="may_19_release"
	release.Type="new"
	release.Owner="roopa@gmail.com"
	release.Project="9"
	release.TargetDate=time.Now()
	log.Print("release", release)

	createdRelease := db.Debug().Create(&release)
	var errMessage = createdRelease.Error
	log.Print("error release", errMessage)

	if createdRelease.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "release created successful": &release})
}

func GetAllReleases (c *gin.Context)  {
	release := models.Release{}
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

	for rows.Next() {
		release:= &models.Release{}
		err := db.Debug().ScanRows(rows, &release)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", release)
		releaseArr=append(releaseArr, release)
	}


	//for _,releaseData := range releaseArr{
	//	fmt.Printf("%v\n", releaseData)
	//}
	if rows.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "releases": releaseArr})
}
