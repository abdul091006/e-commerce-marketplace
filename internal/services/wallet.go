package services

import (
	"fmt"
	"strconv"

	"e-commerce_marketplace/internal/models"
	"e-commerce_marketplace/internal/repositories"
	"e-commerce_marketplace/pkg/utils"
)

type WalletService interface {
	CreateWallet(walletUserID string) (*models.Wallet, error)
	GetWallet(walletUserID string) (*models.Wallet, error)
	AddBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error)
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
	exists, err := s.walletRepo.ExistsByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewWalletError(utils.CodeWalletExists, "Wallet already exists for this user", "")
	}

	// get all type from frappe
	types, err := utils.GetAllBalanceTypesFromFrappe()
	if err != nil {
		return nil, err
	}

	var initialBalances models.BalanceData
	if initialBalances == nil {
		initialBalances = make(models.BalanceData)
	}
	for _, t := range types {
		initialBalances[t] = 0
	}

	// create new wallet
	wallet := &models.Wallet{WalletUserID: walletUserID}
	if err := wallet.SetBalances(&initialBalances); err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to set initial balances", err.Error())
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWallet(walletUserID string) (*models.Wallet, error) {
	return s.walletRepo.GetByWalletUserID(walletUserID)
}

func (s *walletService) AddBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error) {
	// validasi request struct
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, utils.NewWalletError(utils.CodeValidationError, "Validation failed", "")
	}

	// parse amount
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInvalidAmount, "amount must be a number", err.Error())
	}

	// validasi balance type via Frappe
	if err := utils.ValidateBalanceTypeFromFrappe(req.BalanceType); err != nil {
		return nil, err
	}
	if err := utils.ValidateAmount(amountFloat); err != nil {
		return nil, err
	}

	// ambil wallet
	wallet, err := s.walletRepo.GetByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}

	balances, err := wallet.GetBalances()
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to parse balances", err.Error())
	}
	if *balances == nil {
		*balances = make(models.BalanceData)
	}

	if _, ok := (*balances)[req.BalanceType]; !ok {
		(*balances)[req.BalanceType] = 0
	}

	// add balance
	(*balances)[req.BalanceType] += amountFloat

	// update DB
	if err := s.walletRepo.UpdateBalances(walletUserID, balances); err != nil {
		return nil, err
	}

	return s.walletRepo.GetByWalletUserID(walletUserID)
}

func (s *walletService) DeductBalance(walletUserID string, req *utils.UpdateBalanceRequest) (*models.Wallet, error) {
	// validasi request struct
	if validationErrors := utils.ValidateStruct(req); len(validationErrors) > 0 {
		return nil, utils.NewWalletError(utils.CodeValidationError, "Validation failed", "")
	}

	// parse amount
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInvalidAmount, "amount must be a number", err.Error())
	}

	// validasi balance type via Frappe
	if err := utils.ValidateBalanceTypeFromFrappe(req.BalanceType); err != nil {
		return nil, err
	}
	if err := utils.ValidateAmount(amountFloat); err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.GetByWalletUserID(walletUserID)
	if err != nil {
		return nil, err
	}

	// get balance
	balances, err := wallet.GetBalances()
	if err != nil {
		return nil, utils.NewWalletError(utils.CodeInternalError, "Failed to parse balances", err.Error())
	}
	if *balances == nil {
		*balances = make(models.BalanceData)
	}

	if _, ok := (*balances)[req.BalanceType]; !ok {
		(*balances)[req.BalanceType] = 0
	}

	// cek sufficient balance
	if (*balances)[req.BalanceType] < amountFloat {
		return nil, utils.NewWalletError(
			utils.CodeInsufficientBalance,
			fmt.Sprintf("Insufficient %s balance", req.BalanceType),
			fmt.Sprintf("Current balance: %.2f, required: %.2f", (*balances)[req.BalanceType], amountFloat),
		)
	}

	// deduct balance
	(*balances)[req.BalanceType] -= amountFloat

	// update DB
	if err := s.walletRepo.UpdateBalances(walletUserID, balances); err != nil {
		return nil, err
	}

	return s.walletRepo.GetByWalletUserID(walletUserID)
}
