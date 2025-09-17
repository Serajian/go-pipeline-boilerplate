package middleware

import (
	"context"

	"go-pipeline/config"
	"go-pipeline/pkg/generate"

	"github.com/gin-gonic/gin"
)

func TraceIDGenerator() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := generate.TraceID()

		currentCTX := c.Request.Context()
		newCTX := context.WithValue(currentCTX, config.TraceIDKey, traceID)

		c.Request = c.Request.WithContext(newCTX)

		c.Next()
	}
}
