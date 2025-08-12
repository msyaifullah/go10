// main.go
package main

import (
	"flag"
	"log"

	"loan-service/internal/application"
	"loan-service/internal/routes"
	"loan-service/pkg/config"
	"loan-service/pkg/database"
	"loan-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	env := flag.String("env", "local", "Environment (local, staging, production)")
	configPath := flag.String("config", "./configs", "Config file path")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath, *env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger with masking
	appLogger := logger.NewLogger(cfg.Logger)

	// Initialize database connection
	db, err := database.NewConnection(cfg.Database, appLogger)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer db.Close()

	// Initialize application
	app := application.NewApplication().
		WithConfig(cfg).
		WithLogger(appLogger).
		WithDatabase(db).
		WithRepositories().
		WithAdapters().
		WithServices().
		WithHandlers()

	// Validate application
	if err := app.Validate(); err != nil {
		appLogger.Fatal("Application validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Setup Gin
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Configure trusted proxies to eliminate security warning
	if len(cfg.Server.TrustedProxies) > 0 {
		r.SetTrustedProxies(cfg.Server.TrustedProxies)
	} else {
		// If no trusted proxies configured, trust only localhost for development
		r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	r.Use(logger.GinMiddleware(appLogger))
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(r, app)

	// Start cron service
	app.CronService.Start()
	defer app.CronService.Stop()

	appLogger.Info("Server starting", map[string]interface{}{
		"port":        cfg.Server.Port,
		"environment": cfg.App.Environment,
	})

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		appLogger.Fatal("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
