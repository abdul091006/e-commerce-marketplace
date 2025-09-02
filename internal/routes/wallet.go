package routes

import (
	"e-commerce_marketplace/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func WalletRoutes(app *fiber.App, walletHandler *handlers.WalletHandler) {
	// API version prefix
	api := app.Group("/api/v1")
	
	// Wallet routes
	wallets := api.Group("/wallets")
	
	// POST /api/v1/wallets - Create a new wallet
	wallets.Post("/", walletHandler.CreateWallet)
	
	// GET /api/v1/wallets/:id - Get wallet by ID
	wallets.Get("/:id", walletHandler.GetWallet)
	
	// POST /api/v1/wallets/:id/add - Add balance
	wallets.Post("/:id/add", walletHandler.AddBalance)
	
	// POST /api/v1/wallets/:id/deduct - Deduct balance
	wallets.Post("/:id/deduct", walletHandler.DeductBalance)
}