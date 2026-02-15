package middleware

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Log request
        end := time.Now()
        latency := end.Sub(start)
        
        if raw != "" {
            path = path + "?" + raw
        }
        
        log.Printf("[%s] %s %s %d %v", c.Request.Method, path, c.ClientIP(), c.Writer.Status(), latency)
    }
}