package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/interfaces/contract" // For service contracts
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/pkg/http_response" // For error handling and response utilities
	"net/http"                          // Import net/http for HTTP status codes

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// RouterConfig holds dependencies for the router.
type RouterConfig struct {
	AuthService       contract.AuthService
	ProjectService    contract.ProjectService
	StationService    contract.StationService
	LithologyService  contract.LithologyService
	LaboratoryService contract.LaboratoryService
	ReportService     contract.ReportService
	JWTSecret         string
}

// SetupRouter initializes and configures the Gin router.
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	router := gin.Default()

	// Global Middleware
	router.Use(gin.Logger())   // Logging requests
	router.Use(gin.Recovery()) // Catches panics and recovers

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(cfg.AuthService)
	projectHandler := handler.NewProjectHandler(cfg.ProjectService)
	stationHandler := handler.NewStationHandler(cfg.StationService)
	lithologyHandler := handler.NewLithologyHandler(cfg.LithologyService)
	laboratoryHandler := handler.NewLaboratoryHandler(cfg.LaboratoryService)
	reportHandler := handler.NewReportHandler(cfg.ReportService)

	// Initialize Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.AuthService)

	// API Grouping
	api := router.Group("/api")
	{
		// Setup routes for each domain using dedicated functions
		setupAuthRoutes(api, authHandler, authMiddleware)
		setupProjectRoutes(api, projectHandler, stationHandler, authMiddleware)
		// Pass all relevant handlers to setupStationRoutes for nested resources
		setupStationRoutes(api, stationHandler, lithologyHandler, laboratoryHandler, authMiddleware)
		// These now only handle general top-level resources, nested ones are in setupStationRoutes/setupLabSampleRoutes
		setupLithologyRoutes(api, lithologyHandler, authMiddleware)
		setupLabSampleRoutes(api, laboratoryHandler, authMiddleware)
		setupUCSResultRoutes(api, laboratoryHandler, authMiddleware)
		setupReportRoutes(api, reportHandler, authMiddleware)
	}

	// Basic health check route
	router.GET("/health", func(c *gin.Context) {
		http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

// setupAuthRoutes configures authentication-related API endpoints.
func setupAuthRoutes(apiGroup *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		// Authenticated routes
		authGroup.Use(authMiddleware.Handle())
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/profile", authHandler.GetProfile)
		authGroup.PUT("/profile", authHandler.UpdateProfile)
	}
}

// setupProjectRoutes configures project-related API endpoints.
func setupProjectRoutes(apiGroup *gin.RouterGroup, projectHandler *handler.ProjectHandler, stationHandler *handler.StationHandler, authMiddleware *middleware.AuthMiddleware) {
	projectGroup := apiGroup.Group("/projects")
	projectGroup.Use(authMiddleware.Handle()) // Apply auth middleware to all project routes
	{
		// IMPORTANT: Register more specific routes (with more segments) before general wildcard routes.
		// Standardized to use ':id' for project ID consistently.
		projectGroup.GET("/:id/stations", stationHandler.ListStationsByProject) // Specific route for stations under a project

		projectGroup.POST("", projectHandler.CreateProject)
		projectGroup.GET("", projectHandler.ListProjects)   // List all projects for user
		projectGroup.GET("/:id", projectHandler.GetProject) // General route for getting a single project by ID
		projectGroup.PUT("/:id", projectHandler.UpdateProject)
		projectGroup.DELETE("/:id", projectHandler.DeleteProject)
	}
}

// setupStationRoutes configures general station-related API endpoints and its nested resources.
// It now defines routes for lithology logs and lab samples that belong to a specific station.
func setupStationRoutes(apiGroup *gin.RouterGroup, stationHandler *handler.StationHandler, lithologyHandler *handler.LithologyHandler, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware) {
	stationGroup := apiGroup.Group("/stations")
	stationGroup.Use(authMiddleware.Handle()) // Apply auth middleware
	{
		stationGroup.POST("", stationHandler.CreateStation)
		stationGroup.GET("/:id", stationHandler.GetStation)
		stationGroup.PUT("/:id", stationHandler.UpdateStation)
		stationGroup.DELETE("/:id", stationHandler.DeleteStation)

		// Station-specific Lithology Log routes (using :id for station)
		stationGroup.GET("/:id/lithology-logs", lithologyHandler.ListLithologyLogsByStation)
		stationGroup.GET("/:id/lithology-logs/by-depth", lithologyHandler.ListLithologyLogsByDepthRange)

		// Station-specific Lab Sample routes (using :id for station)
		stationGroup.GET("/:id/lab-samples", laboratoryHandler.ListLabSamplesByStation)
	}
}

// setupLithologyRoutes configures lithology log-related API endpoints (general, top-level).
// Station-specific lithology routes are now handled in setupStationRoutes.
func setupLithologyRoutes(apiGroup *gin.RouterGroup, lithologyHandler *handler.LithologyHandler, authMiddleware *middleware.AuthMiddleware) {
	lithologyGroup := apiGroup.Group("/lithology-logs")
	lithologyGroup.Use(authMiddleware.Handle())
	{
		lithologyGroup.POST("", lithologyHandler.CreateLithologyLog)
		lithologyGroup.GET("/:id", lithologyHandler.GetLithologyLogByID)
		lithologyGroup.PUT("/:id", lithologyHandler.UpdateLithologyLog)
		lithologyGroup.DELETE("/:id", lithologyHandler.DeleteLithologyLog)
	}
}

// setupLabSampleRoutes configures laboratory sample-related API endpoints (general, top-level)
// and its nested UCS results. Station-specific lab sample routes are now in setupStationRoutes.
func setupLabSampleRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware) {
	labSampleGroup := apiGroup.Group("/lab-samples")
	labSampleGroup.Use(authMiddleware.Handle())
	{
		labSampleGroup.POST("", laboratoryHandler.CreateLabSample)
		labSampleGroup.GET("/:id", laboratoryHandler.GetLabSampleByID)
		labSampleGroup.PUT("/:id", laboratoryHandler.UpdateLabSample)
		labSampleGroup.DELETE("/:id", laboratoryHandler.DeleteLabSample)

		// Specific route for UCS results under a lab sample (using :id for sample)
		labSampleGroup.GET("/:id/ucs-results", laboratoryHandler.ListUCSResultsBySample)
	}
}

// setupUCSResultRoutes configures UCS result-related API endpoints (general, top-level).
// Sample-specific UCS result routes are now handled in setupLabSampleRoutes.
func setupUCSResultRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware) {
	ucsResultGroup := apiGroup.Group("/ucs-results")
	ucsResultGroup.Use(authMiddleware.Handle())
	{
		ucsResultGroup.POST("", laboratoryHandler.CreateUCSResult)
		ucsResultGroup.GET("/:id", laboratoryHandler.GetUCSResultByID)
		ucsResultGroup.PUT("/:id", laboratoryHandler.UpdateUCSResult)
		ucsResultGroup.DELETE("/:id", laboratoryHandler.DeleteUCSResult)
	}
}

// setupReportRoutes configures report generation API endpoints.
func setupReportRoutes(apiGroup *gin.RouterGroup, reportHandler *handler.ReportHandler, authMiddleware *middleware.AuthMiddleware) {
	reportGroup := apiGroup.Group("/reports")
	reportGroup.Use(authMiddleware.Handle())
	{
		reportGroup.GET("/stations/:station_id", reportHandler.GenerateStationReport)
		reportGroup.GET("/projects/:id/summary", reportHandler.GenerateProjectSummary) // Standardized to :id for project ID
	}
}
