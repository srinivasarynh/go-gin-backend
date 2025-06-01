package middleware

import (
    "time"
    "go-gin-backend/pkg/logger"

    "github.com/gin-gonic/gin"
)

func Logger(logger *logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()

        if raw != "" {
            path = path + "?" + raw
        }

        logger.Info("HTTP Request",
            map[string]interface{}{
                "method":     method,
                "path":       path,
                "status":     statusCode,
                "latency":    latency,
                "client_ip":  clientIP,
                "user_agent": c.Request.UserAgent(),
            })
    }
}

func Recovery() gin.HandlerFunc {
    return gin.Recovery()
}
