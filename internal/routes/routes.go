package routes

import (
    "go-gin-backend/internal/handlers"
    "go-gin-backend/internal/middleware"
    "go-gin-backend/internal/services"

    "github.com/gin-gonic/gin"
)

func SetupRoutes(
    r *gin.Engine,
    authHandler *handlers.AuthHandler,
    userHandler *handlers.UserHandler,
    healthHandler *handlers.HealthHandler,
    authService *services.AuthService,
) {
    // Health check routes
    r.GET("/health", healthHandler.Health)
    r.GET("/ready", healthHandler.Ready)

    // API routes
    api := r.Group("/api/v1")
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
        }

        // Protected routes
        protected := api.Group("/")
        protected.Use(middleware.AuthMiddleware(authService))
        {
            // Auth protected routes
            protected.GET("/auth/me", authHandler.Me)

            // User routes
            users := protected.Group("/users")
            {
                users.GET("", userHandler.ListUsers)
                users.GET("/:id", userHandler.GetUser)
                users.PUT("/:id", userHandler.UpdateUser)
                users.DELETE("/:id", userHandler.DeleteUser)
            }
        }
    }
}
