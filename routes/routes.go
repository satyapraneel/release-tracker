package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/controllers"
)

// Config will hold repositories that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R *gin.Engine
}

//RouterGin function
func RouterGin(app *cmd.Application) *gin.Engine {
	log.Println("****name*******: ", app.Name)
	releaseHandler := controllers.NewReleaseHandler(app)
	router := gin.Default()
	router.Static("/assets", "./ui/assets")
	router.LoadHTMLGlob("ui/html/**/*.tmpl")
	// router.LoadHTMLGlob("ui/html/*.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home", gin.H{
			"title": "Release tracker",
		})
	})
	api := router.Group("/release")
	{
		api.GET("/list", releaseHandler.GetListOfReleases)
		api.GET("/create", releaseHandler.CreateReleaseForm)
		api.POST("/store", releaseHandler.CreateRelease)
		api.GET("/getReviewers", releaseHandler.GetProjectReviewerList)
		//api.GET("/users/:id", user.GetUser)
		//api.PUT("/users/:id", user.UpdateUser)
		//api.DELETE("/users/:id", user.DeleteUser)
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	return router
}
