package route

import (
	"asyvy/web/internal/controller"
	"net/http"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var db = make(map[string]string)

func InitRoute(r *gin.Engine, c *controller.Controller) {
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.POST("/imagescan", c.ImageScan)
}
