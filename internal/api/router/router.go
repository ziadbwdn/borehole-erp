package router

import (
	"boreholedata-ms/internal/api/handler"
	"boreholedata-ms/internal/interfaces/contract"
	"boreholedata-ms/internal/logger" // Import logger
	"boreholedata-ms/internal/middleware"
	"boreholedata-ms/internal/models" // Import models for UserRole constants
	"boreholedata-ms/pkg/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RouterConfig holds dependencies for the router.
type RouterConfig struct {
	AuthService         contract.AuthService
	ProjectService      contract.ProjectService
	StationService      contract.StationService
	LithologyService    contract.LithologyService
	LaboratoryService   contract.LaboratoryService
	ReportService       contract.ReportService
	UserActivityService contract.UserActivityService // <--- ADDED: UserActivityService
	JWTSecret           string
	Logger              logger.Logger // Add Logger to RouterConfig
}

// SetupRouter initializes and configures the Gin router.
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	router := gin.Default()

	// Global Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize Handlers
	authHandler := handler.NewAuthHandler(cfg.AuthService)
	projectHandler := handler.NewProjectHandler(cfg.ProjectService)
	stationHandler := handler.NewStationHandler(cfg.StationService)
	lithologyHandler := handler.NewLithologyHandler(cfg.LithologyService)
	laboratoryHandler := handler.NewLaboratoryHandler(cfg.LaboratoryService)
	reportHandler := handler.NewReportHandler(cfg.ReportService)
	userActivityHandler := handler.NewUserActivityHandler(cfg.UserActivityService, cfg.Logger) // <--- ADDED: UserActivityHandler initialization

	// Initialize Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.AuthService)

	// API Grouping
	api := router.Group("/api")
	{
		// Setup routes for each domain using dedicated functions, passing the logger for authorization middleware
		setupAuthRoutes(api, authHandler, authMiddleware, cfg.Logger)
		setupProjectRoutes(api, projectHandler, stationHandler, authMiddleware, cfg.Logger)
		setupStationRoutes(api, stationHandler, lithologyHandler, laboratoryHandler, authMiddleware, cfg.Logger)
		setupLithologyRoutes(api, lithologyHandler, authMiddleware, cfg.Logger)
		setupLabSampleRoutes(api, laboratoryHandler, authMiddleware, cfg.Logger)
		setupUCSResultRoutes(api, laboratoryHandler, authMiddleware, cfg.Logger)
		setupReportRoutes(api, reportHandler, authMiddleware, cfg.Logger)
		setupUserActivityRoutes(api, userActivityHandler, authMiddleware, cfg.Logger) // <--- ADDED: UserActivity routes setup
	}

	// Basic health check route
	router.GET("/health", func(c *gin.Context) {
		http_response.RespondWithSuccess(c, http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

// setupAuthRoutes configures authentication-related API endpoints.
// No specific role-based authorization needed here, as it's just general auth ops.
func setupAuthRoutes(apiGroup *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		// Authenticated routes
		authGroup.Use(authMiddleware.Handle()) // Apply authentication middleware
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/profile", authHandler.GetProfile)
		authGroup.PUT("/profile", authHandler.UpdateProfile)
	}
}

// setupProjectRoutes configures project-related API endpoints.
func setupProjectRoutes(apiGroup *gin.RouterGroup, projectHandler *handler.ProjectHandler, stationHandler *handler.StationHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	projectGroup := apiGroup.Group("/projects")
	projectGroup.Use(authMiddleware.Handle()) // All project routes require authentication
	{
		// Project Status Updates: Exclusive to Engineers (or Project Managers).
		// Assuming UpdateProject might include status updates, we apply the role check here.
		projectGroup.POST("", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.CreateProject) // Assuming only engineers can create
		projectGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.UpdateProject)
		projectGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleEngineer), projectHandler.DeleteProject) // Assuming only engineers can delete

		// General read access can be for all authenticated users (Engineer, Geologist, Lab Technician)
		projectGroup.GET("", projectHandler.ListProjects)
		projectGroup.GET("/:id", projectHandler.GetProject)
		projectGroup.GET("/:id/stations", stationHandler.ListStationsByProject) // List stations under a project
	}
}

// setupStationRoutes configures general station-related API endpoints and its nested resources.
func setupStationRoutes(apiGroup *gin.RouterGroup, stationHandler *handler.StationHandler, lithologyHandler *handler.LithologyHandler, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	stationGroup := apiGroup.Group("/stations")
	stationGroup.Use(authMiddleware.Handle()) // All station routes require authentication
	{
		// Drilling Status (Station) Updates: Handled by both Engineers and Geologists.
		stationGroup.POST("", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.CreateStation) // Assuming creation also requires these roles
		stationGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.UpdateStation)
		stationGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleEngineer, models.RoleGeologist), stationHandler.DeleteStation)

		// General read access for stations
		stationGroup.GET("/:id", stationHandler.GetStation)

		// Station-specific Lithology Log routes (using :id for station)
		// Lithology Log Updates: Solely managed by Geologists.
		stationGroup.GET("/:id/lithology-logs", lithologyHandler.ListLithologyLogsByStation)
		stationGroup.GET("/:id/lithology-logs/by-depth", lithologyHandler.ListLithologyLogsByDepthRange) // Read access

		// Station-specific Lab Sample routes (using :id for station)
		// Laboratory Sample Updates: Responsibility of Lab Technicians.
		stationGroup.GET("/:id/lab-samples", laboratoryHandler.ListLabSamplesByStation) // Read access
	}
}

// setupLithologyRoutes configures lithology log-related API endpoints (general, top-level).
func setupLithologyRoutes(apiGroup *gin.RouterGroup, lithologyHandler *handler.LithologyHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	lithologyGroup := apiGroup.Group("/lithology-logs")
	lithologyGroup.Use(authMiddleware.Handle()) // All lithology log routes require authentication
	{
		// Lithology Log Updates: Solely managed by Geologists.
		lithologyGroup.POST("", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.CreateLithologyLog)
		lithologyGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.UpdateLithologyLog)
		lithologyGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleGeologist), lithologyHandler.DeleteLithologyLog)

		// General read access
		lithologyGroup.GET("/:id", lithologyHandler.GetLithologyLogByID)
	}
}

// setupLabSampleRoutes configures laboratory sample-related API endpoints (general, top-level)
func setupLabSampleRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	labSampleGroup := apiGroup.Group("/lab-samples")
	labSampleGroup.Use(authMiddleware.Handle()) // All lab sample routes require authentication
	{
		// Laboratory Sample Updates: Responsibility of Lab Technicians.
		labSampleGroup.POST("", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.CreateLabSample)
		labSampleGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.UpdateLabSample)
		labSampleGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.DeleteLabSample)

		// General read access
		labSampleGroup.GET("/:id", laboratoryHandler.GetLabSampleByID)

		// Specific route for UCS results under a lab sample (using :id for sample) - read access
		labSampleGroup.GET("/:id/ucs-results", laboratoryHandler.ListUCSResultsBySample)
	}
}

// setupUCSResultRoutes configures UCS result-related API endpoints (general, top-level).
func setupUCSResultRoutes(apiGroup *gin.RouterGroup, laboratoryHandler *handler.LaboratoryHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	ucsResultGroup := apiGroup.Group("/ucs-results")
	ucsResultGroup.Use(authMiddleware.Handle()) // All UCS result routes require authentication
	{
		// Assuming UCS results are part of Lab Technician's responsibility.
		ucsResultGroup.POST("", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.CreateUCSResult)
		ucsResultGroup.PUT("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.UpdateUCSResult)
		ucsResultGroup.DELETE("/:id", middleware.AuthorizeRole(log, models.RoleLabTechnician), laboratoryHandler.DeleteUCSResult)

		// General read access
		ucsResultGroup.GET("/:id", laboratoryHandler.GetUCSResultByID)
	}
}

// setupReportRoutes configures report generation API endpoints.
func setupReportRoutes(apiGroup *gin.RouterGroup, reportHandler *handler.ReportHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	reportGroup := apiGroup.Group("/reports")
	reportGroup.Use(authMiddleware.Handle()) // All report routes require authentication
	{
		// Reports can generally be accessed by all roles for viewing.
		// If specific reports are sensitive, you'd add AuthorizeRole here.
		reportGroup.GET("/stations/:station_id", reportHandler.GenerateStationReport)
		reportGroup.GET("/projects/:id/summary", reportHandler.GenerateProjectSummary)
	}
}

// New: setupUserActivityRoutes configures user activity related API endpoints.
func setupUserActivityRoutes(apiGroup *gin.RouterGroup, userActivityHandler *handler.UserActivityHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	activityGroup := apiGroup.Group("/activities")
	activityGroup.Use(authMiddleware.Handle()) // All activity routes require authentication
	{
		// LogActivity: Can be called by any authenticated user for their own activities, or by admin/engineer for audit.
		// For simplicity, we'll allow all authenticated users to log their own activities.
		// If you want stricter control (e.g., only internal services log), adjust roles.
		activityGroup.POST("", userActivityHandler.LogActivity) // No specific role check here for general logging

		// ListActivities, GetActivitySummary, GetSecurityAlerts, CleanOldActivities:
		// These are typically for auditing/admin purposes.
		// Assuming only Engineers (or an Admin role if you had one) can view and manage activity logs.
		activityGroup.GET("", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.ListActivities)
		activityGroup.GET("/summary/:userID", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.GetActivitySummary)
		activityGroup.GET("/alerts", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.GetSecurityAlerts)
		activityGroup.DELETE("", middleware.AuthorizeRole(log, models.RoleAdmin), userActivityHandler.CleanOldActivities) // Admin-level cleanup
	}
}
