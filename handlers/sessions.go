package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/config"
)

func SetupSession(router *gin.Engine) {
	store := cookie.NewStore([]byte(config.SessionDetails().Secret))
	router.Use(sessions.Sessions("mysession", store))
}
