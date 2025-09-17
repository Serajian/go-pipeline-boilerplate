package http

import (
	"go-pipeline/config"
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
	"go-pipeline/internal/presentation/http/middleware"

	"github.com/gin-gonic/gin"
)

type GinAdapter struct {
	Engin       *gin.Engine
	pipeline    ports.Pipeline[model.UserData]
	shortRunner ports.PipeLineFn[model.UserData]
	barrier     ports.PipeLineBarrier[model.UserData]
}

func NewGinAdapter(
	p ports.Pipeline[model.UserData],
	b ports.PipeLineBarrier[model.UserData],
	sr ports.PipeLineFn[model.UserData],
) *GinAdapter {
	adapter := &GinAdapter{
		Engin:       ginEngin(),
		pipeline:    p,
		barrier:     b,
		shortRunner: sr,
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

	g.testParallel(layer)
	g.testBarrier(layer)
	g.testShort(layer)
}

func selectMode(debug bool) string {
	if debug {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

func (g *GinAdapter) testParallel(r *gin.RouterGroup) {
	r.GET("/v1", func(c *gin.Context) {
		ctx := c.Request.Context()
		user := model.UserData{
			Name: "mohsenV1",
			Age:  30,
			// Email: "hooora!",
			Email: "V1@gmailcom",
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

func (g *GinAdapter) testBarrier(r *gin.RouterGroup) {
	r.GET("/v2", func(c *gin.Context) {
		ctx := c.Request.Context()
		user := model.UserData{
			Name:  "mohsenV2",
			Age:   30,
			Email: "V2@gmailcom",
		}

		in := make(chan model.UserData, 1)
		in <- user
		close(in)
		out, errChan := g.barrier.Run(ctx, in)

		var result []model.UserData

		for out != nil || errChan != nil {
			select {
			case <-ctx.Done():
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
				c.JSON(500, gin.H{"error": e.Error()})
				_ = e

			}
		}
		c.JSON(200, gin.H{"data": result, "count": len(result)})
	})
}

func (g *GinAdapter) testShort(r *gin.RouterGroup) {
	r.GET("/v3", func(c *gin.Context) {
		ctx := c.Request.Context()
		user := model.UserData{
			Name: "mohsenV3",
			Age:  30,
			// Email: "hooora!",
			Email: "V3@gmailcom",
		}

		out, err := g.shortRunner.Run(ctx, user)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": out})
	})
}
