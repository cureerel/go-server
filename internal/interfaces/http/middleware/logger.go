package middleware

import (
	"github.com/cureerel/cserver/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		requestID := c.GetString("request_id")

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		fields := []logger.Field{
			{Key: "request_id", Value: requestID},
			{Key: "path", Value: path},
			{Key: "method", Value: c.Request.Method},
			{Key: "status", Value: status},
			{Key: "duration_ms", Value: duration.Milliseconds()},
			{Key: "ip", Value: c.ClientIP()},
			{Key: "user_agent", Value: c.Request.UserAgent()},
		}

		if len(c.Errors) > 0 {
			fields = append(fields, logger.Field{Key: "errors", Value: c.Errors})
			log.Error("HTTP Request Failed", fields...)
			return
		}

		log.Info("HTTP Request", fields...)
	}
}
