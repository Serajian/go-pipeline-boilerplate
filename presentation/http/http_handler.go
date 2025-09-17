package http

import (
	"go-pipeline/config"

	"github.com/gin-gonic/gin"
)

type GinAdapter struct {
	Engin *gin.Engine
}

func NewGinAdapter() *GinAdapter {
	adapter := &GinAdapter{
		Engin: ginEngin(),
	}
	adapter.handleRoutes()
	return adapter
}

func ginEngin() *gin.Engine {
	gin.SetMode(selectMode(config.Get().AppConfig.Debug))

	router := gin.New()

	if config.Get().AppConfig.Debug {
		router.Use(gin.Logger())
	}

	return router
}

func (g *GinAdapter) handleRoutes() {
	// TODO:1: change name to yours

	layer := g.Engin.Group("/boiler")
	layer.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})
}

func selectMode(debug bool) string {
	if debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
