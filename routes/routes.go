package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/cmd"
	"github.com/release-trackers/gin/controllers"
	"github.com/release-trackers/gin/handlers"
	"github.com/release-trackers/gin/middleware"
)

// Config will hold repositories that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R *gin.Engine
}

//RouterGin function
func RouterGin(app *cmd.Application) {
	log.Println("****name*******: ", app.Name)
	handler := controllers.NewHandler(app)
	router := gin.Default()
	handlers.SetupSession(router)
	router.Static("/assets", "./ui/assets")
	router.LoadHTMLGlob("ui/html/**/*.tmpl")
	// router.LoadHTMLGlob("ui/html/*.tmpl")

	router.GET("/login", handler.LoginForm)
	router.Use(middleware.ParseForm()).POST("/login", handler.Login)
	router.GET("/logout", controllers.Logout)
	auth := router.Group("/")

	auth.Use(middleware.Authentication())
	{
		auth.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home", gin.H{
				"title": "Release tracker",
			})
		})
		api := auth.Group("/release")
		{
			api.GET("/index", handler.GetIndex)
			api.POST("/list", handler.GetListOfReleases)
			api.GET("/create", handler.CreateReleaseForm)
			api.POST("/store", handler.CreateRelease)
			api.GET("/getReviewers", handler.GetProjectReviewerList)
			api.GET("/show/:id", handler.ViewReleaseForm)
			//api.PUT("/users/:id", user.UpdateUser)
			//api.DELETE("/users/:id", user.DeleteUser)
		}
	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.Run(os.Getenv("PORT"))
}
