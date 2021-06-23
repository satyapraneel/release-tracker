package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/release-trackers/gin/config"
)

func SetupSession(router *gin.Engine) {
	store := cookie.NewStore([]byte(config.SessionDetails().Secret))
	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
}
