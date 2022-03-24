package api

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter(e *gin.Engine) {
	v1 := e.Group("/v1")
	e.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1.POST("/compress", Compress)
}
