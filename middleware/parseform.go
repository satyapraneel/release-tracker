package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		if err != nil {
			http.Error(c.Writer, "Bad Request", http.StatusBadRequest)
			return
		}
	}
}
