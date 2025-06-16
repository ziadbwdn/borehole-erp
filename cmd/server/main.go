package main

import (
	"boreholedata-ms/config"
	"boreholedata-ms/internal/api/router"
	"boreholedata-ms/internal/database"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/usecases/repository" // This will now refer to gorm_user_activity_repository
	"boreholedata-ms/internal/usecases/services"

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
	// --- 1. Load Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	// --- 1.5. Initialize Application Logger ---
	appLogger := logger.New()
	appLogger.Info(context.Background(), "Application logger initialized.")

	// --- 2. Initialize Database Connection ---
	// InitDB is from your /internal/database/connection.go
	db, appErr := database.InitDB(cfg)
	if appErr != nil {
		appLogger.Error(context.Background(), "Failed to connect to database", appErr)
		log.Fatalf("Failed to connect to database: %v", appErr.Error())
	}
	appLogger.Info(context.Background(), "Database connection established.")

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Station{},
		&models.LithologyLog{},
		&models.LabSample{},
		&models.UCSResult{},
		&models.UserActivity{}, // Ensure UserActivity model is included for migration
	)
	if err != nil {
		appLogger.Error(context.Background(), "Failed to auto-migrate database", err)
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	appLogger.Info(context.Background(), "Database auto-migration completed (if enabled).")

	// --- 3. Initialize Repositories ---
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	stationRepo := repository.NewStationRepository(db)
	lithologyRepo := repository.NewLithologyRepository(db)
	labRepo := repository.NewLaboratoryRepository(db)
	userActivityRepo := repository.NewGormUserActivityRepository(db, appLogger) // <--- CHANGED HERE

	// --- 4. Initialize Services ---
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	projectService := services.NewProjectService(projectRepo)
	stationService := services.NewStationService(stationRepo, projectRepo)
	lithologyService := services.NewLithologyService(lithologyRepo, stationRepo)
	laboratoryService := services.NewLaboratoryService(labRepo, stationRepo)
	reportService := services.NewReportService(stationRepo, lithologyRepo, labRepo)
	userActivityService := services.NewUserActivityService(userActivityRepo, appLogger)

	// --- 5. Setup Router ---
	routerConfig := &router.RouterConfig{
		AuthService:         authService,
		ProjectService:      projectService,
		StationService:      stationService,
		LithologyService:    lithologyService,
		LaboratoryService:   laboratoryService,
		ReportService:       reportService,
		UserActivityService: userActivityService,
		JWTSecret:           cfg.JWTSecret,
		Logger:              appLogger,
	}
	r := router.SetupRouter(routerConfig)

	// --- 6. Start HTTP Server ---
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
	appLogger.Info(context.Background(), "Access API at /api/*")
	appLogger.Info(context.Background(), "Swagger UI documentation available at /swagger/index.html")

	// --- 7. Graceful Shutdown ---
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
