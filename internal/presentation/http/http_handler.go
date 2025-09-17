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
	pipeline ports.Pipeline[model.UserData]
}

func NewGinAdapter(pipeline ports.Pipeline[model.UserData]) *GinAdapter {
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
			// Email: "hooora@gmailcom",
		}
		in := make(chan model.UserData, 1)
		in <- user
		close(in)

		out, errChan := g.pipeline.Chain(ctx, in)

		var result []model.UserData
		done := ctx.Done()

		for out != nil || errChan != nil {
			select {
			case <-done:
				c.JSON(499, gin.H{"error": "client canceled"})
				return
			case m, ok := <-out:
				if !ok {
					out = nil
					continue
				}
				result = append(result, m)
			case e, ok := <-errChan:
				if !ok {
					errChan = nil
					continue
				}
				// TODO: you can collect for log
				c.JSON(500, gin.H{"error": e.Error()})
				_ = e
				return
			}
		}
		c.JSON(200, gin.H{"data": result, "count": len(result)})
	})
}

func selectMode(debug bool) string {
	if debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
