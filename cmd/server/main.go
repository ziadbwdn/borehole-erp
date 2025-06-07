package main

import (
	"boreholedata-ms/config" // Import the new root-level config package
	"boreholedata-ms/internal/api/router"
	"boreholedata-ms/internal/database" // Import the database connection package
	"boreholedata-ms/internal/models"
	"boreholedata-ms/internal/usecases/repository"
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

	// --- 2. Initialize Database Connection ---
	// Use the InitDB function from the database package
	db, appErr := database.InitDB(cfg) // Pass the loaded config
	if appErr != nil {
		log.Fatalf("Failed to connect to database: %v", appErr.Error()) // Use appErr.Error() for custom error
	}
	log.Println("Database connection established.")

	// Auto-migrate models (Optional, for development or initial setup)
	// This will create or update tables based on your GORM struct definitions.
	// In production, consider using a proper migration tool.
	err = db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Station{},
		&models.LithologyLog{},
		&models.LabSample{},
		&models.UCSResult{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	log.Println("Database auto-migration completed (if enabled).")

	// --- 3. Initialize Repositories ---
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	stationRepo := repository.NewStationRepository(db)
	lithologyRepo := repository.NewLithologyRepository(db)
	labRepo := repository.NewLaboratoryRepository(db)

	// --- 4. Initialize Services ---
	authService := services.NewAuthService(userRepo, cfg.JWTSecret) // Use JWTSecret from config
	projectService := services.NewProjectService(projectRepo)
	stationService := services.NewStationService(stationRepo, projectRepo)
	lithologyService := services.NewLithologyService(lithologyRepo, stationRepo)
	laboratoryService := services.NewLaboratoryService(labRepo, stationRepo)
	reportService := services.NewReportService(stationRepo, lithologyRepo, labRepo)

	// --- 5. Setup Router ---
	routerConfig := &router.RouterConfig{
		AuthService:       authService,
		ProjectService:    projectService,
		StationService:    stationService,
		LithologyService:  lithologyService,
		LaboratoryService: laboratoryService,
		ReportService:     reportService,
		JWTSecret:         cfg.JWTSecret, // Use JWTSecret from config
	}
	r := router.SetupRouter(routerConfig)

	// --- 6. Start HTTP Server ---
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port), // Use Port from config
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to listen: %s\n", err)
		}
	}()

	log.Printf("Server is running on port %s", cfg.Port)
	log.Println("Access API at /api/*")
	log.Println("Swagger UI documentation available at /swagger/index.html")

	// --- 7. Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown due to an error: %v", err)
	}

	log.Println("Server exited successfully.")
}
