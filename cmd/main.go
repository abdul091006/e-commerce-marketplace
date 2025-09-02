package main

import (
	"log"
	"os"

	"e-commerce_marketplace/internal/config"
	"e-commerce_marketplace/internal/handlers"
	"e-commerce_marketplace/internal/repositories"
	"e-commerce_marketplace/internal/routes"
	"e-commerce_marketplace/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	walletRepo := repositories.NewWalletRepository(db)

	// Initialize services
	walletService := services.NewWalletService(walletRepo)

	// Initialize handlers
	walletHandler := handlers.NewWalletHandler(walletService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	routes.WalletRoutes(app, walletHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}