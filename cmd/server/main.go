// @title Borehole Data Management API
// @version 1.0
// @description This is the Borehole Data Management API documentation.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Example: "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."

// @host localhost:8080
// @BasePath /api
// @schemes http

package main

import (
	"boreholedata-ms/config"
	"boreholedata-ms/internal/api/router"
	"boreholedata-ms/internal/database"
	"boreholedata-ms/internal/logger"
	_ "boreholedata-ms/docs"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// --- 1. Load Config ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- 1.5. Initialize App Logger ---
	appLogger := logger.New()
	appLogger.Info(context.Background(), "Application logger initialized.")

	// --- 2. Initialize DB Connection ---
	db, appErr := database.InitDB(cfg)
	if appErr != nil {
		appLogger.Error(context.Background(), "Failed to connect to database", appErr)
		log.Fatalf("Failed to connect to database: %v", appErr.Error())
	}
	appLogger.Info(context.Background(), "Database connection established.")

	// --- 2.5. Run DB Migrations ---
	if err := database.RunMigrations(db); err != nil {
		appLogger.Error(context.Background(), "Failed to run database migrations", err)
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	appLogger.Info(context.Background(), "Database migrations completed successfully.")

	// --- 3. Setup Router (Now with DI) --- responsible for creating its own dependencies.
	routerConfig := &router.RouterConfig{
		DB:     db,
		Cfg:    cfg,
		Logger: appLogger,
	}
	r := router.SetupRouter(routerConfig)
	appLogger.Info(context.Background(), "Router and dependencies initialized.")

	// --- 4. Start Server ---
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error(context.Background(), "Server failed to listen", err)
			log.Fatalf("Server failed to listen: %s\n", err)
		}
	}()

	appLogger.Info(context.Background(), fmt.Sprintf("Server is running on port %s", cfg.Port))
	appLogger.Info(context.Background(), "Swagger UI documentation available at /swagger/index.html")

	// --- 5. Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info(context.Background(), "Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error(context.Background(), "Server forced to shutdown due to an error", err)
		log.Fatalf("Server forced to shutdown due to an error: %v", err)
	}

	appLogger.Info(context.Background(), "Server exited successfully.")
}
