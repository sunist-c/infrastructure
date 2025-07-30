package plugins

import (
	"context"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/gin-gonic/gin"
)

func GinTraceRequest(c *gin.Context) {
	c.Request = c.Request.WithContext(
		context.WithValue(
			trace.Trace(c.Request.Context()),
			trace.RequestType,
			&trace.Request{
				Method:   c.Request.Method,
				Path:     c.Request.URL.Path,
				Host:     c.Request.URL.Host,
				ClientIP: c.ClientIP(),
			},
		),
	)
}
