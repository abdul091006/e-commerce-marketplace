package repositories

import (
	"e-commerce_marketplace/internal/models"
	"e-commerce_marketplace/pkg/utils"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

type WalletRepository interface {
	// Create creates a new wallet
	Create(wallet *models.Wallet) error
	
	// GetByWalletUserID retrieves a wallet by wallet user ID
	GetByWalletUserID(walletUserID string) (*models.Wallet, error)
	
	// Update updates an existing wallet
	Update(wallet *models.Wallet) error
	
	// Delete soft deletes a wallet by wallet user ID
	Delete(walletUserID string) error
	
	// ExistsByWalletUserID checks if a wallet exists for the given wallet user ID
	ExistsByWalletUserID(walletUserID string) (bool, error)
	
	// UpdateBalances updates the balances field of a wallet
	UpdateBalances(walletUserID string, balances *models.BalanceData) error
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		if isUniqueConstraintError(err) {
			return utils.NewWalletError(utils.CodeWalletExists, "Wallet already exists for this user", err.Error())
		}
		return utils.NewWalletError(utils.CodeDatabaseError, "Failed to create wallet", err.Error())
	}
	return nil
}

func (r *walletRepository) GetByWalletUserID(walletUserID string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.First(&wallet, "wallet_user_id = ?", walletUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.NewWalletError(utils.CodeWalletNotFound, "Wallet not found", "")
		}
		return nil, utils.NewWalletError(utils.CodeDatabaseError, "Failed to retrieve wallet", err.Error())
	}
	return &wallet, nil
}

func (r *walletRepository) Update(wallet *models.Wallet) error {
	if err := r.db.Save(wallet).Error; err != nil {
		if isUniqueConstraintError(err) {
			return utils.NewWalletError(utils.CodeWalletExists, "Wallet already exists for this user", err.Error())
		}
		return utils.NewWalletError(utils.CodeDatabaseError, "Failed to update wallet", err.Error())
	}
	return nil
}

func (r *walletRepository) Delete(walletUserID string) error {
	result := r.db.Delete(&models.Wallet{}, "wallet_user_id = ?", walletUserID)
	if result.Error != nil {
		return utils.NewWalletError(utils.CodeDatabaseError, "Failed to delete wallet", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return utils.NewWalletError(utils.CodeWalletNotFound, "Wallet not found", "")
	}
	return nil
}

func (r *walletRepository) ExistsByWalletUserID(walletUserID string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Wallet{}).Where("wallet_user_id = ?", walletUserID).Count(&count).Error; err != nil {
		return false, utils.NewWalletError(utils.CodeDatabaseError, "Failed to check wallet existence", err.Error())
	}
	return count > 0, nil
}

func (r *walletRepository) UpdateBalances(walletUserID string, balances *models.BalanceData) error {
	// Convert balances to JSON
	balancesData, err := json.Marshal(balances)
	if err != nil {
		return utils.NewWalletError(utils.CodeInternalError, "Failed to marshal balances", err.Error())
	}

	// Update only the balances field
	result := r.db.Model(&models.Wallet{}).
		Where("wallet_user_id = ?", walletUserID).
		Update("balances", datatypes.JSON(balancesData))
	if result.Error != nil {
		return utils.NewWalletError(utils.CodeDatabaseError, "Failed to update wallet balances", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return utils.NewWalletError(utils.CodeWalletNotFound, "Wallet not found", "")
	}
	return nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	// PostgreSQL unique constraint error contains "duplicate key value"
	return err != nil && 
		(containsString(err.Error(), "duplicate key value") || 
		 containsString(err.Error(), "UNIQUE constraint failed"))
}

// containsString checks if a string contains a substring (case-insensitive)
func containsString(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    str[:len(substr)] == substr || 
		    str[len(str)-len(substr):] == substr ||
		    containsSubstring(str, substr))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}