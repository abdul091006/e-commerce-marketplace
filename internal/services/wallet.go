package services

import (
	"e-commerce_marketplace/internal/models"
	"e-commerce_marketplace/internal/repositories"
	"e-commerce_marketplace/pkg/utils"
	"fmt"

	"strconv"
)

type WalletService interface {
	// CreateWallet creates a new wallet for a user
	CreateWallet(walletUserID string) (*models.Wallet, error)
	
	// GetWallet retrieves a wallet by wallet user ID
	GetWallet(walletUserID string) (*models.Wallet, error)
	
	// AddBalance adds amount to specified balance type
	AddBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error)
	
	// DeductBalance deducts amount from specified balance type
	DeductBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error)
}

type walletService struct {
	walletRepo repositories.WalletRepository
}

func NewWalletService(walletRepo repositories.WalletRepository) WalletService {
	return &walletService{
		walletRepo: walletRepo,
	}
}

func (s *walletService) CreateWallet(walletUserID string) (*models.Wallet, error) {
	// Cek apakah wallet sudah ada untuk user ini
	exists, err := s.walletRepo.ExistsByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewWalletError(utils.CodeWalletExists, "Wallet already exists for this user", "")
	}

	// Buat initial balances
	initialBalances := &models.BalanceData{
		Coins: 0,
		EXP:   0,
	}

	// Buat wallet baru
	wallet := &models.Wallet{
		WalletUserID: walletUserID,
	}

	// Set balances awal
	if err := wallet.SetBalances(initialBalances); err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to set initial balances", err.Error())
	}

	// Simpan ke repository
	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *walletService) GetWallet(walletUserID string) (*models.Wallet, error) {
	return s.walletRepo.GetByWalletUserID(walletUserID)
}

func (s *walletService) AddBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, utils.NewWalletError(utils.CodeValidationError, "Validation failed", "")
	}

	amountFloat, _ := strconv.ParseFloat(req.Amount, 64)

	// Additional validation
	if err := utils.ValidateBalanceType(req.BalanceType); err != nil {
		return nil, err
	}
	if err := utils.ValidateAmount(amountFloat); err != nil {
		return nil, err
	}

	// Get existing wallet
	wallet, err := s.walletRepo.GetByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}

	// Get current balances
	balances, err := wallet.GetBalances()
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to parse current balances", err.Error())
	}

	// Update balance based on type
	switch req.BalanceType {
	case "coins":
		balances.Coins += amountFloat
	case "exp":
		balances.EXP += amountFloat
	}

	// Update wallet balances
	if err := s.walletRepo.UpdateBalances(walletUserID, balances); err != nil {
		return nil, err
	}

	// Return updated wallet
	return s.walletRepo.GetByWalletUserID(walletUserID)
}

func (s *walletService) DeductBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error) {
	// Validate request
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, utils.NewWalletError(utils.CodeValidationError, "Validation failed", "")
	}

	amountFloat, _ := strconv.ParseFloat(req.Amount, 64)

	// Additional validation
	if err := utils.ValidateBalanceType(req.BalanceType); err != nil {
		return nil, err
	}
	if err := utils.ValidateAmount(amountFloat); err != nil {
		return nil, err
	}

	// Get existing wallet
	wallet, err := s.walletRepo.GetByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}

	// Get current balances
	balances, err := wallet.GetBalances()
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to parse current balances", err.Error())
	}

	// Check and update balance based on type
	switch req.BalanceType {
	case "coins":
		if balances.Coins < amountFloat {
			return nil, utils.NewWalletError(
				utils.CodeInsufficientBalance,
				"Insufficient coin balance",
				fmt.Sprintf("Current balance: %.2f, required: %.2f", balances.Coins, amountFloat),
			)
		}
		balances.Coins -= amountFloat
	case "exp":
		if balances.EXP < amountFloat {
			return nil, utils.NewWalletError(
				utils.CodeInsufficientBalance,
				"Insufficient EXP balance",
				fmt.Sprintf("Current balance: %.2f, required: %.2f", balances.EXP, amountFloat),
			)
		}
		balances.EXP -= amountFloat
	}

	// Update wallet balances
	if err := s.walletRepo.UpdateBalances(walletUserID, balances); err != nil {
		return nil, err
	}

	// Return updated wallet
	return s.walletRepo.GetByWalletUserID(walletUserID)
}