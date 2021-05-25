package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/controllers/release"
)

// Config will hold repositories that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R *gin.Engine
}

//RouterGin function
func RouterGin() *gin.Engine {

	router := gin.Default()
	router.Static("/assets", "./ui/assets")
	router.LoadHTMLGlob("ui/html/**/*.tmpl")
	// router.LoadHTMLGlob("ui/html/*.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home/index", gin.H{
			"title": "Main website",
		})
	})
	api := router.Group("/release")
	{
		api.GET("/list", release.GetListOfReleases)
		api.GET("/create", release.CreateReleaseForm)
		api.POST("/store", release.CreateRelease)
		//api.GET("/users/:id", user.GetUser)
		//api.PUT("/users/:id", user.UpdateUser)
		//api.DELETE("/users/:id", user.DeleteUser)
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	return router
}
