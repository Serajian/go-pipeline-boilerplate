package http

import (
	"go-pipeline/config"
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
	"go-pipeline/internal/presentation/http/middleware"

	"github.com/gin-gonic/gin"
)

type GinAdapter struct {
	Engin    *gin.Engine
	pipeline ports.RegistryPipeline
}

func NewGinAdapter(pipeline ports.RegistryPipeline) *GinAdapter {
	adapter := &GinAdapter{
		Engin:    ginEngin(),
		pipeline: pipeline,
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
	layer.Use(middleware.TraceIDGenerator())
	layer.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		user := model.UserData{
			Name:  "mohsen",
			Age:   30,
			Email: "hooora!",
		}
		in := make(chan model.UserData, 1)
		in <- user
		if err := g.pipeline.Run(ctx, in); err != nil {
			c.JSON(200, gin.H{"msg": err.Error()})
		}
		c.JSON(200, gin.H{"msg": "hello world"})
	})
}

func selectMode(debug bool) string {
	if debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
