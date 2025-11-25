package main

import (
	"log"

	"github.com/brunobarlari/inventorypulse/internal/config"
	"github.com/brunobarlari/inventorypulse/internal/handler"
	"github.com/brunobarlari/inventorypulse/internal/middleware"
	"github.com/brunobarlari/inventorypulse/internal/repository"
	"github.com/brunobarlari/inventorypulse/internal/service"
	"github.com/brunobarlari/inventorypulse/pkg/database"
	"github.com/brunobarlari/inventorypulse/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// @title           InventoryPulse API
// @version         1.0
// @description     Backend API for inventory management with real-time updates via WebSocket
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@inventorypulse.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connection
	db, err := database.NewPostgresConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get underlying SQL DB for defer close
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Run seeder
	if err := database.RunSeeder(db, cfg); err != nil {
		log.Fatalf("Failed to run seeder: %v", err)
	}

	// Initialize JWT service
	jwtService := jwt.NewJWTService(&cfg.JWT)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "InventoryPulse API is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to InventoryPulse API v1.0",
			})
		})

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)

			// Protected routes (require authentication)
			authProtected := auth.Group("")
			authProtected.Use(authMiddleware.RequireAuth())
			{
				authProtected.GET("/me", authHandler.Me)

				// Admin only routes
				authAdmin := authProtected.Group("")
				authAdmin.Use(authMiddleware.RequireAdmin())
				{
					authAdmin.POST("/register", authHandler.Register)
				}
			}
		}

		// Protected API routes will be added here
		// Categories and Products routes will go here in next phases
	}

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
