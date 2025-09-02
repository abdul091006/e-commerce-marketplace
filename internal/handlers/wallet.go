package handlers

import (
	"e-commerce_marketplace/pkg/utils"
	"e-commerce_marketplace/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService services.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// CreateWallet handles POST /wallets
func (h *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	userID := uuid.New().String()

	wallet, err := h.walletService.CreateWallet(userID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return utils.CreatedResponse(c, "Wallet created successfully", wallet)
}

// GetWallet handles GET /wallets/:id
func (h *WalletHandler) GetWallet(c *fiber.Ctx) error {
	walletUserID := c.Params("id")
	if walletUserID == "" {
		return utils.BadRequestResponse(c, "Wallet user ID is required", "")
	}

	wallet, err := h.walletService.GetWallet(walletUserID)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, "Wallet retrieved successfully", wallet)
}

// AddBalance handles POST /wallets/:id/add
func (h *WalletHandler) AddBalance(c *fiber.Ctx) error {
	walletUserID := c.Params("id")
	if walletUserID == "" {
		return utils.BadRequestResponse(c, "Wallet user ID is required", "")
	}

	var req utils.UpdateBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	wallet, err := h.walletService.AddBalance(walletUserID, &req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, "Balance added successfully", wallet)
}

// DeductBalance handles POST /wallets/:id/deduct
func (h *WalletHandler) DeductBalance(c *fiber.Ctx) error {
	walletUserID := c.Params("id")
	if walletUserID == "" {
		return utils.BadRequestResponse(c, "Wallet user ID is required", "")
	}

	var req utils.UpdateBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	wallet, err := h.walletService.DeductBalance(walletUserID, &req)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, "Balance deducted successfully", wallet)
}

// handleServiceError converts service errors to appropriate HTTP responses
func (h *WalletHandler) handleServiceError(c *fiber.Ctx, err error) error {
	if utils.IsWalletError(err) {
		walletErr := err.(*utils.WalletError)
		switch walletErr.Code {
		case utils.CodeWalletNotFound:
			return utils.NotFoundResponse(c, walletErr.Message)
		case utils.CodeWalletExists:
			return utils.ConflictResponse(c, walletErr.Message)
		case utils.CodeInsufficientBalance:
			return utils.BadRequestResponse(c, walletErr.Message, walletErr.Details)
		case utils.CodeInvalidAmount, utils.CodeInvalidBalanceType, utils.CodeValidationError:
			return utils.BadRequestResponse(c, walletErr.Message, walletErr.Details)
		default:
			return utils.InternalServerErrorResponse(c, "An error occurred while processing your request")
		}
	}
	return utils.InternalServerErrorResponse(c, "An unexpected error occurred")
}