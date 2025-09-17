package http

import (
	"net/http"

	"go-pipeline/config"
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
	"go-pipeline/internal/presentation/http/middleware"

	"github.com/gin-gonic/gin"
)

type GinAdapter struct {
	Engin       *gin.Engine
	pipeline    ports.ChainPipeline[model.UserData]
	shortRunner ports.ShortCircuitPipeLine[model.UserData]
	barrier     ports.BarrierPipeLine[model.UserData]
}

func NewGinAdapter(
	p ports.ChainPipeline[model.UserData],
	b ports.BarrierPipeLine[model.UserData],
	sr ports.ShortCircuitPipeLine[model.UserData],
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

		items, errs, canceled := drainAll(ctx, out, errChan)
		if canceled {
			c.JSON(499, gin.H{"error": "client canceled"})
			return
		}

		if errs != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errs})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"pipeline": "parallel",
			"count":    len(items),
			"items":    items,
		})
	})
}

func (g *GinAdapter) testBarrier(r *gin.RouterGroup) {
	r.GET("/v2", func(c *gin.Context) {
		ctx := c.Request.Context()

		in := make(chan model.UserData, 1)
		in <- model.UserData{Name: "mohsenV2", Age: 30, Email: "V2@gmailcom"}
		close(in)

		out, errChan := g.barrier.Run(ctx, in)

		items, errs, canceled := drainAll(ctx, out, errChan)
		if canceled {
			c.JSON(499, gin.H{"error": "client canceled"})
			return
		}
		if errs != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errs})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"pipeline": "barrier",
			"count":    len(items),
			"items":    items,
		})
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
