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
	"github.com/brunobarlari/inventorypulse/pkg/websocket"
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

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	wsHandler := websocket.NewHandler(wsHub)

	// Initialize JWT service
	jwtService := jwt.NewJWTService(&cfg.JWT)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtService)
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, wsHub)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":           "ok",
			"message":          "InventoryPulse API is running",
			"websocket_clients": wsHub.GetClientCount(),
		})
	})

	// WebSocket endpoint
	router.GET("/ws", wsHandler.HandleWebSocket)

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

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(authMiddleware.RequireAuth())
			{
				authProtected.GET("/me", authHandler.Me)

				// Admin only
				authAdmin := authProtected.Group("")
				authAdmin.Use(authMiddleware.RequireAdmin())
				{
					authAdmin.POST("/register", authHandler.Register)
				}
			}
		}

		// Category routes
		categories := api.Group("/categories")
		categories.Use(authMiddleware.RequireAuth())
		{
			categories.GET("", categoryHandler.List)
			categories.GET("/:id", categoryHandler.Get)

			// Admin only
			categoriesAdmin := categories.Group("")
			categoriesAdmin.Use(authMiddleware.RequireAdmin())
			{
				categoriesAdmin.POST("", categoryHandler.Create)
				categoriesAdmin.PUT("/:id", categoryHandler.Update)
				categoriesAdmin.DELETE("/:id", categoryHandler.Delete)
			}
		}

		// Product routes
		products := api.Group("/products")
		products.Use(authMiddleware.RequireAuth())
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Get)

			// Admin only
			productsAdmin := products.Group("")
			productsAdmin.Use(authMiddleware.RequireAdmin())
			{
				productsAdmin.POST("", productHandler.Create)
				productsAdmin.PUT("/:id", productHandler.Update)
				productsAdmin.DELETE("/:id", productHandler.Delete)
				productsAdmin.PATCH("/:id/stock", productHandler.UpdateStock)
			}
		}
	}

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("WebSocket available at ws://localhost:%s/ws", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
