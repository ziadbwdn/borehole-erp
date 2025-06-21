package router

import (
	"boreholedata-ms/config"
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/logger"
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/usecases/repository"
	"boreholedata-ms/internal/usecases/services"
	"boreholedata-ms/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// RouterConfig holds the core dependencies needed to build the API layer.
type RouterConfig struct {
	DB     *gorm.DB
	Cfg    *config.Config
	Logger logger.Logger
}

// SetupRouter initializes and configures the Gin router.
// It is the central point for dependency injection and route setup orchestration.
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	// --- 1. Initialize Dependencies ---
	userRepo := repository.NewUserRepository(cfg.DB)
	projectRepo := repository.NewProjectRepository(cfg.DB)
	stationRepo := repository.NewStationRepository(cfg.DB)
	lithologyRepo := repository.NewLithologyRepository(cfg.DB)
	labRepo := repository.NewLaboratoryRepository(cfg.DB)
	userActivityRepo := repository.NewGormUserActivityRepository(cfg.DB, cfg.Logger)

	userActivityService := services.NewUserActivityService(userActivityRepo, cfg.Logger)
	authService := services.NewAuthService(userRepo, userActivityService, cfg.Cfg.JWTSecret)
	projectService := services.NewProjectService(projectRepo, userActivityService, cfg.Logger)
	stationService := services.NewStationService(stationRepo, projectRepo, userActivityService, cfg.Logger)
	lithologyService := services.NewLithologyService(lithologyRepo, stationRepo, userActivityService, cfg.Logger)
	laboratoryService := services.NewLaboratoryService(labRepo, stationRepo, userActivityService, cfg.Logger)
	reportService := services.NewReportService(stationRepo, lithologyRepo, labRepo, userActivityService, cfg.Logger)

	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService, authService)
	stationHandler := handler.NewStationHandler(stationService, authService)
	lithologyHandler := handler.NewLithologyHandler(lithologyService, authService)
	laboratoryHandler := handler.NewLaboratoryHandler(laboratoryService, authService)
	reportHandler := handler.NewReportHandler(reportService, authService)
	userActivityHandler := handler.NewUserActivityHandler(userActivityService, cfg.Logger)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	// --- 2. Setup Router Engine ---
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"status": "ok"})
	})

	// --- 3. Setup API Route Groups ---
	api := router.Group("/api")
	{
		// Public routes are set up directly
		setupAuthRoutes(api.Group("/auth"), authHandler, authMiddleware)

		// Protected routes are grouped under the auth middleware
		protectedAPI := api.Group("/")
		protectedAPI.Use(authMiddleware.Handle())
		{
			// Call the setup functions for each domain
			setupProjectRoutes(protectedAPI, projectHandler, stationHandler, cfg.Logger)
			setupStationRoutes(protectedAPI, stationHandler, lithologyHandler, laboratoryHandler, cfg.Logger)
			setupLithologyRoutes(protectedAPI, lithologyHandler, cfg.Logger)
			setupLabSampleRoutes(protectedAPI, laboratoryHandler, cfg.Logger)
			setupUCSResultRoutes(protectedAPI, laboratoryHandler, cfg.Logger)
			setupReportRoutes(protectedAPI, reportHandler, cfg.Logger)
			setupUserActivityRoutes(protectedAPI, userActivityHandler, cfg.Logger)
		}
	}

	return router
}