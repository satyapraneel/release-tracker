package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
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
	handlers.SetupSession(router, "release_tracker")
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
			session := sessions.Default(c)
			flashes := session.Flashes()
			session.Save()
			c.HTML(http.StatusOK, "home", gin.H{
				"title":   "Release tracker",
				"flashes": flashes,
			})
		})
	}
	api := auth.Group("/release")
	{
		api.GET("/index", handler.GetIndex)
		api.POST("/list", handler.GetListOfReleases)
		api.GET("/create", handler.CreateReleaseForm)
		api.POST("/store", handler.CreateRelease)
		api.GET("/getReviewers", handler.GetProjectReviewerList)
		api.GET("/show/:id", handler.ViewReleaseForm)
		api.GET("/tickets", handler.ReleaseTicketsForm)
		api.GET("/getTickets", handler.ReleaseListTicketsByReleaseId)
	}

	projects := auth.Group("/projects")
	{
		projects.GET("/", handler.GetProjects)
		projects.POST("/list", handler.GetListOfProjects)
		projects.GET("/create", handler.CreateProjectForm)
		projects.GET("/store", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/projects/create")
		})
		projects.POST("/store", handler.CreateProject)
		projects.GET("/show/:id", handler.ViewProjectForm)
		projects.POST("/update/:id", handler.UpdateProject)
		projects.GET("/update/:id", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/projects")
		})
	}

	reviewers := auth.Group("/reviewers")
	{
		reviewers.GET("/", handler.GetReviewers)
		reviewers.POST("/list", handler.GetListOfReviewers)
		reviewers.GET("/create", handler.CreateReviewersForm)
		reviewers.GET("/store", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/reviewers/create")
		})
		reviewers.POST("/store", handler.CreateReviewer)
		reviewers.GET("/show/:id", handler.ViewReviewerForm)
		reviewers.POST("/update/:id", handler.UpdateReviewer)
		reviewers.GET("/delete/:id", handler.DeleteReviewer)
		reviewers.GET("/update/:id", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/reviewers")
		})
	}

	dls := auth.Group("/dls")
	{
		dls.GET("/", handler.GetDLs)
		dls.POST("/list", handler.GetListOfDLs)
		dls.GET("/create", handler.CreateDLsForm)
		dls.GET("/store", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/dls/create")
		})
		dls.POST("/store", handler.CreateDL)
		dls.GET("/show/:id", handler.ViewDLForm)
		dls.POST("/update/:id", handler.UpdateDL)
		dls.GET("/delete/:id", handler.DeleteDL)
		dls.GET("/update/:id", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/dls")
		})
	}


	//oauthapi := auth.Group("/oauth")
	//{
	//	oauthapi.GET("/index", handler.GetAuthCode)
	//	oauthapi.POST("/code", handler.GetAccessToken)
	//}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.Run(os.Getenv("PORT"))
}
