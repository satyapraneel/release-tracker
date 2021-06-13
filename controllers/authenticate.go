package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/repositories"
	"github.com/release-trackers/gin/validation"
)

func (app *App) Login(c *gin.Context) {

	validation.ValidateLoginForm(c)
	session := sessions.Default(c)
	repositories := repositories.NewRepositoryHandler(app.Application)
	user, err := repositories.AuthenticateUser(c.Request.PostForm.Get("email"), c.Request.PostForm.Get("password"))
	if err {
		c.HTML(http.StatusUnauthorized, "auth/login", gin.H{"error": "Unauthorized"})
		c.Abort()
	}
	session.Set("id", user.ID)
	session.Set("email", user.Email)
	session.Save()
	c.Redirect(http.StatusFound, "/")

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func (app *App) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/login", gin.H{})
}
