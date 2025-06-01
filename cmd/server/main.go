package main

import (
	"context"
	"go-gin-backend/internal/config"
	"go-gin-backend/internal/database"
	"go-gin-backend/internal/handlers"
	"go-gin-backend/internal/middleware"
	"go-gin-backend/internal/repository"
	"go-gin-backend/internal/routes"
	"go-gin-backend/internal/services"
	"go-gin-backend/pkg/logger"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	appLogger := logger.NewLogger()
	defer appLogger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatal("failed to load config", err)
	}

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("failed to connect to database", err)
	}

	if err := database.RunMigrations(cft.DatabaseURL); err != nil {
		appLogger.Fatal("failed to run migrations", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Logger(appLogger))
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())

	routes.SetupRoutes(r, authHandler, userHandler, healthHandler, authService)

	srv := &http.Server{
		Addr : ":" + cfg.Port,
		Handler: r,
	}

go func() {
        appLogger.Info("Starting server on port " + cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.Fatal("Failed to start server", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    appLogger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30* time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        appLogger.Fatal("Server forced to shutdown", err)
    }

    appLogger.Info("Server exited")
}
