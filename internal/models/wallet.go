package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Wallet struct {
	WalletUserID string         `json:"wallet_user_id" gorm:"unique;not null;index"`
	Balances     datatypes.JSON `json:"balances" gorm:"type:jsonb;default:'{\"coins\":0,\"exp\":0}'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Wallet model
func (Wallet) TableName() string {
	return "wallets"
}

// BalanceData represents the structure of the balances JSON field
type BalanceData struct {
	Coins float64 `json:"coins"`
	EXP   float64 `json:"exp"`
}

// GetBalances parses the JSONB balances field into BalanceData
func (w *Wallet) GetBalances() (*BalanceData, error) {
	var balances BalanceData
	if err := json.Unmarshal(w.Balances, &balances); err != nil {
		return nil, err
	}
	return &balances, nil
}

// SetBalances sets the balances field from BalanceData
func (w *Wallet) SetBalances(balances *BalanceData) error {
	data, err := json.Marshal(balances)
	if err != nil {
		return err
	}
	w.Balances = datatypes.JSON(data)
	return nil
}
