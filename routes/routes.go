package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/controllers/release"
)

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R *gin.Engine
}

//RouterGin function
func RouterGin() *gin.Engine {

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})
	api := router.Group("/api")
	{
		api.GET("/releases", release.GetAllReleases)
		api.POST("/release/create", release.CreateUser)
		//api.GET("/users/:id", user.GetUser)
		//api.PUT("/users/:id", user.UpdateUser)
		//api.DELETE("/users/:id", user.DeleteUser)
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	return router
}
