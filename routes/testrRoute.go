package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetGetPong(r *gin.Engine) {
	test := r.Group("/")
	test.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ponger")
	 })
}