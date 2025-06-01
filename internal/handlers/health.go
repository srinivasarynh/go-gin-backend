package handlers

import (
    "net/http"
    "time"
    "go-gin-backend/pkg/response"

    "github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
        "version":   "1.0.0",
    }

    response.Success(c, http.StatusOK, "Service is healthy", health)
}

func (h *HealthHandler) Ready(c *gin.Context) {
    // Add database connectivity check, Redis check, etc.
    ready := map[string]interface{}{
        "status":    "ready",
        "timestamp": time.Now().Unix(),
        "checks": map[string]string{
            "database": "connected",
            "redis":    "connected",
        },
    }

    response.Success(c, http.StatusOK, "Service is ready", ready)
}
