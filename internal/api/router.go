package api

import (
	"net/http"
	"portfolio-be/internal/api/handlers"
	"portfolio-be/internal/api/middleware"
	"portfolio-be/internal/config"
	"portfolio-be/internal/repository"
	"portfolio-be/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, s3Service *services.S3Service, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Portfolio Backend API is running",
		})
	})

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize repositories
	contentRepo := repository.NewContentRepository(db)
	uploadRepo := repository.NewUploadRepository(db)
	resourceRepo := repository.NewResourceRepository(db)
	experienceRepo := repository.NewExperienceRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	technologyRepo := repository.NewTechnologyRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	testimonialRepo := repository.NewTestimonialRepository(db)
	userRepo := repository.NewUserRepository(db)
	contactRepo := repository.NewContactRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// Initialize services
	contentService := services.NewContentService(contentRepo)
	uploadService := services.NewUploadService(uploadRepo, s3Service)
	resourceService := services.NewResourceService(resourceRepo, uploadRepo, s3Service)
	experienceService := services.NewExperienceService(experienceRepo)
	serviceService := services.NewServiceService(serviceRepo)
	technologyService := services.NewTechnologyService(technologyRepo)
	projectService := services.NewProjectService(projectRepo)
	testimonialService := services.NewTestimonialService(testimonialRepo)
	jwtService := services.NewJWTService(cfg.JWTConfig.SecretKey, cfg.JWTConfig.Issuer)
	authService := services.NewAuthService(userRepo, jwtService)
	contactService := services.NewContactService(contactRepo)
	roleService := services.NewRoleService(roleRepo, permissionRepo, userRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	// Initialize middleware
	permissionMiddleware := middleware.NewPermissionMiddleware(userRepo)

	// Initialize Cron Service
	cronService := services.NewCronService(resourceService, uploadService)
	// Start cron service in background
	go cronService.Start()

	// Initialize handlers
	contentHandler := handlers.NewContentHandler(contentService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	resourceHandler := handlers.NewResourceHandler(resourceService)
	experienceHandler := handlers.NewExperienceHandler(experienceService)
	serviceHandler := handlers.NewServiceHandler(serviceService)
	technologyHandler := handlers.NewTechnologyHandler(technologyService)
	projectHandler := handlers.NewProjectHandler(projectService)
	testimonialHandler := handlers.NewTestimonialHandler(testimonialService)
	portfolioHandler := handlers.NewPortfolioHandler(experienceService, serviceService, technologyService, projectService, testimonialService)
	authHandler := handlers.NewAuthHandler(authService, permissionMiddleware)
	userHandler := handlers.NewUserHandler(userRepo)
	contactHandler := handlers.NewContactHandler(contactService)
	statsHandler := handlers.NewStatsHandler(projectService, experienceService, technologyService, serviceService, testimonialService, contactService)
	adminOrderHandler := handlers.NewAdminOrderHandler(projectService, experienceService, technologyService, serviceService, testimonialService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Portfolio endpoint (combined data)
	router.GET("/api/portfolio", portfolioHandler.GetPortfolioData)

	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Protected auth routes
		authProtected := auth.Group("")
		authProtected.Use(middleware.AuthMiddleware(jwtService))
		{
			authProtected.GET("/profile", authHandler.Profile)
			authProtected.POST("/logout", authHandler.Logout)
		}
	}

	// Admin routes (protected)
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtService))
	admin.Use(middleware.AdminMiddleware())
	{
		// User management
		admin.GET("/users", permissionMiddleware.RequirePermission("users", "read"), userHandler.GetUsers)
		admin.POST("/users", permissionMiddleware.RequirePermission("users", "create"), userHandler.CreateUser)
		admin.GET("/users/:id", permissionMiddleware.RequirePermission("users", "read"), userHandler.GetUser)
		admin.PUT("/users/:id", permissionMiddleware.RequirePermission("users", "update"), userHandler.UpdateUser)
		admin.PATCH("/users/:id/password", permissionMiddleware.RequirePermission("users", "update"), userHandler.UpdateUserPassword)
		admin.DELETE("/users/:id", permissionMiddleware.RequirePermission("users", "delete"), userHandler.DeleteUser)
		admin.PATCH("/users/:id/toggle-status", permissionMiddleware.RequirePermission("users", "update"), userHandler.ToggleUserStatus)
		admin.POST("/users/assign-role", permissionMiddleware.RequirePermission("users", "update"), userHandler.AssignRole)
		admin.GET("/users/:id/permissions", permissionMiddleware.RequirePermission("users", "read"), userHandler.GetUserPermissions)

		// Role management
		admin.GET("/roles", permissionMiddleware.RequirePermission("roles", "read"), roleHandler.GetAllRoles)
		admin.POST("/roles", permissionMiddleware.RequirePermission("roles", "create"), roleHandler.CreateRole)
		admin.GET("/roles/:id", permissionMiddleware.RequirePermission("roles", "read"), roleHandler.GetRole)
		admin.PUT("/roles/:id", permissionMiddleware.RequirePermission("roles", "update"), roleHandler.UpdateRole)
		admin.DELETE("/roles/:id", permissionMiddleware.RequirePermission("roles", "delete"), roleHandler.DeleteRole)
		admin.POST("/roles/:id/permissions", permissionMiddleware.RequirePermission("roles", "update"), roleHandler.AssignPermissions)
		admin.GET("/roles/:id/permissions", permissionMiddleware.RequirePermission("roles", "read"), roleHandler.GetRolePermissions)

		// Permission management
		admin.GET("/permissions", permissionMiddleware.RequirePermission("permissions", "read"), permissionHandler.GetAllPermissions)
		admin.POST("/permissions", permissionMiddleware.RequirePermission("permissions", "create"), permissionHandler.CreatePermission)
		admin.GET("/permissions/:id", permissionMiddleware.RequirePermission("permissions", "read"), permissionHandler.GetPermission)
		admin.PUT("/permissions/:id", permissionMiddleware.RequirePermission("permissions", "update"), permissionHandler.UpdatePermission)
		admin.DELETE("/permissions/:id", permissionMiddleware.RequirePermission("permissions", "delete"), permissionHandler.DeletePermission)
		admin.GET("/permissions/resource/:resource", permissionMiddleware.RequirePermission("permissions", "read"), permissionHandler.GetPermissionsByResource)
		admin.POST("/permissions/initialize", permissionMiddleware.RequirePermission("permissions", "create"), permissionHandler.InitializeDefaultPermissions)

		// Content management
		admin.POST("/contents", permissionMiddleware.RequirePermission("contents", "create"), contentHandler.CreateContent)
		admin.PUT("/contents/:id", permissionMiddleware.RequirePermission("contents", "update"), contentHandler.UpdateContent)
		admin.DELETE("/contents/:id", permissionMiddleware.RequirePermission("contents", "delete"), contentHandler.DeleteContent)

		// Experience management
		admin.POST("/experiences", permissionMiddleware.RequirePermission("experiences", "create"), experienceHandler.CreateExperience)
		admin.PUT("/experiences/:id", permissionMiddleware.RequirePermission("experiences", "update"), experienceHandler.UpdateExperience)
		admin.DELETE("/experiences/:id", permissionMiddleware.RequirePermission("experiences", "delete"), experienceHandler.DeleteExperience)

		// Service management
		admin.POST("/services", permissionMiddleware.RequirePermission("services", "create"), serviceHandler.CreateService)
		admin.PUT("/services/:id", permissionMiddleware.RequirePermission("services", "update"), serviceHandler.UpdateService)
		admin.DELETE("/services/:id", permissionMiddleware.RequirePermission("services", "delete"), serviceHandler.DeleteService)

		// Technology management
		admin.POST("/technologies", permissionMiddleware.RequirePermission("technologies", "create"), technologyHandler.CreateTechnology)
		admin.PUT("/technologies/:id", permissionMiddleware.RequirePermission("technologies", "update"), technologyHandler.UpdateTechnology)
		admin.DELETE("/technologies/:id", permissionMiddleware.RequirePermission("technologies", "delete"), technologyHandler.DeleteTechnology)

		// Project management
		admin.POST("/projects", permissionMiddleware.RequirePermission("projects", "create"), projectHandler.CreateProject)
		admin.PUT("/projects/:id", permissionMiddleware.RequirePermission("projects", "update"), projectHandler.UpdateProject)
		admin.DELETE("/projects/:id", permissionMiddleware.RequirePermission("projects", "delete"), projectHandler.DeleteProject)

		// Testimonial management
		admin.POST("/testimonials", permissionMiddleware.RequirePermission("testimonials", "create"), testimonialHandler.CreateTestimonial)
		admin.PUT("/testimonials/:id", permissionMiddleware.RequirePermission("testimonials", "update"), testimonialHandler.UpdateTestimonial)
		admin.DELETE("/testimonials/:id", permissionMiddleware.RequirePermission("testimonials", "delete"), testimonialHandler.DeleteTestimonial)

		// Contact management
		admin.GET("/contacts", permissionMiddleware.RequirePermission("contacts", "read"), contactHandler.GetContacts)
		admin.GET("/contacts/:id", permissionMiddleware.RequirePermission("contacts", "read"), contactHandler.GetContact)
		admin.PUT("/contacts/:id", permissionMiddleware.RequirePermission("contacts", "update"), contactHandler.UpdateContact)
		admin.DELETE("/contacts/:id", permissionMiddleware.RequirePermission("contacts", "delete"), contactHandler.DeleteContact)
		admin.GET("/contacts/unread-count", permissionMiddleware.RequirePermission("contacts", "read"), contactHandler.GetUnreadCount)
		admin.PATCH("/contacts/:id/mark-read", permissionMiddleware.RequirePermission("contacts", "update"), contactHandler.MarkAsRead)

		// Upload management
		admin.POST("/uploads", permissionMiddleware.RequirePermission("uploads", "create"), uploadHandler.UploadFile)
		admin.DELETE("/uploads/:id", permissionMiddleware.RequirePermission("uploads", "delete"), uploadHandler.DeleteUpload)

		// Resource management
		admin.POST("/resources", permissionMiddleware.RequirePermission("uploads", "create"), resourceHandler.CreateResource)
		admin.GET("/resources", permissionMiddleware.RequirePermission("uploads", "read"), resourceHandler.GetAllResources)
		admin.GET("/resources/:id", permissionMiddleware.RequirePermission("uploads", "read"), resourceHandler.GetResource)
		admin.PUT("/resources/:id", permissionMiddleware.RequirePermission("uploads", "update"), resourceHandler.UpdateResource)
		admin.DELETE("/resources/:id", permissionMiddleware.RequirePermission("uploads", "delete"), resourceHandler.DeleteResource)
		admin.GET("/resources/stats", permissionMiddleware.RequirePermission("uploads", "read"), resourceHandler.GetResourceStats)
		admin.POST("/resources/refresh-urls", permissionMiddleware.RequirePermission("uploads", "update"), resourceHandler.RefreshExpiredURLs)

		// Order management
		admin.PUT("/projects/order", permissionMiddleware.RequirePermission("projects", "update"), adminOrderHandler.UpdateProjectsOrder)
		admin.PUT("/experiences/order", permissionMiddleware.RequirePermission("experiences", "update"), adminOrderHandler.UpdateExperiencesOrder)
		admin.PUT("/technologies/order", permissionMiddleware.RequirePermission("technologies", "update"), adminOrderHandler.UpdateTechnologiesOrder)
		admin.PUT("/services/order", permissionMiddleware.RequirePermission("services", "update"), adminOrderHandler.UpdateServicesOrder)
		admin.PUT("/testimonials/order", permissionMiddleware.RequirePermission("testimonials", "update"), adminOrderHandler.UpdateTestimonialsOrder)

		// Stats (read-only for most admins)
		admin.GET("/stats", permissionMiddleware.RequireAnyPermission([]string{"projects:read", "experiences:read", "technologies:read", "services:read", "testimonials:read", "contacts:read"}), statsHandler.GetCounts)
	}

	// Public API routes (read-only)
	api := router.Group("/api")
	{
		// Content routes
		api.GET("/contents", contentHandler.GetAllContents)
		api.GET("/contents/:id", contentHandler.GetContent)

		// Experience routes
		api.GET("/experiences", experienceHandler.GetExperiences)
		api.GET("/experiences/:id", experienceHandler.GetExperience)

		// Service routes
		api.GET("/services", serviceHandler.GetAllServices)
		api.GET("/services/:id", serviceHandler.GetService)

		// Technology routes
		api.GET("/technologies", technologyHandler.GetTechnologies)
		api.GET("/technologies/:id", technologyHandler.GetTechnology)

		// Project routes
		api.GET("/projects", projectHandler.GetProjects)
		api.GET("/projects/:id", projectHandler.GetProject)

		// Testimonial routes
		api.GET("/testimonials", testimonialHandler.GetTestimonials)
		api.GET("/testimonials/:id", testimonialHandler.GetTestimonial)

		// Contact routes (for submitting contact forms)
		api.POST("/contacts", contactHandler.CreateContact)

		// Stats routes
		api.GET("/stats/counts", statsHandler.GetCounts)

		// Upload routes
		api.GET("/uploads", uploadHandler.GetAllUploads)
		api.GET("/uploads/:id", uploadHandler.GetUpload)
		api.GET("/uploads/summary", uploadHandler.GetAllUploadsWithSummary)

		// Resource routes (public read-only)
		api.GET("/resources", resourceHandler.GetAllResources)
		api.GET("/resources/:id", resourceHandler.GetResource)
		api.POST("/resources/:id/download", resourceHandler.DownloadResource)
		api.GET("/resources/stats", resourceHandler.GetResourceStats)
	}

	// Legacy API v1 routes (kept for backward compatibility, read-only)
	v1 := router.Group("/api/v1")
	{
		// Content routes
		contents := v1.Group("/contents")
		{
			contents.GET("", contentHandler.GetAllContents)
			contents.GET("/:id", contentHandler.GetContent)
		}

		// Upload routes
		uploads := v1.Group("/uploads")
		{
			uploads.GET("", uploadHandler.GetAllUploads)
			uploads.GET("/:id", uploadHandler.GetUpload)
			uploads.GET("/summary", uploadHandler.GetAllUploadsWithSummary)
		}

		// Resource routes
		resources := v1.Group("/resources")
		{
			resources.GET("", resourceHandler.GetAllResources)
			resources.GET("/:id", resourceHandler.GetResource)
			resources.POST("/:id/download", resourceHandler.DownloadResource)
		}

		// Experience routes
		experiences := v1.Group("/experiences")
		{
			experiences.GET("", experienceHandler.GetExperiences)
			experiences.GET("/:id", experienceHandler.GetExperience)
		}

		// Service routes
		services := v1.Group("/services")
		{
			services.GET("", serviceHandler.GetAllServices)
			services.GET("/:id", serviceHandler.GetService)
		}

		// Technology routes
		technologies := v1.Group("/technologies")
		{
			technologies.GET("", technologyHandler.GetTechnologies)
			technologies.GET("/:id", technologyHandler.GetTechnology)
		}

		// Project routes
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.GET("/:id", projectHandler.GetProject)
		}

		// Testimonial routes
		testimonials := v1.Group("/testimonials")
		{
			testimonials.GET("", testimonialHandler.GetTestimonials)
			testimonials.GET("/:id", testimonialHandler.GetTestimonial)
		}

		// Stats routes
		stats := v1.Group("/stats")
		{
			stats.GET("/counts", statsHandler.GetCounts)
		}
	}

	return router
}
