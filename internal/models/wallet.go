package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Wallet struct {
	WalletUserID string         `json:"wallet_user_id" gorm:"unique;not null;index"`
	Balances     datatypes.JSON `json:"balances"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for the Wallet model
func (Wallet) TableName() string {
	return "wallets"
}

// BalanceData represents the structure of the balances JSON field
type BalanceData map[string]float64


// GetBalances parses the JSON balances field into BalanceData
func (w *Wallet) GetBalances() (*BalanceData, error) {
	var balances BalanceData
	if len(w.Balances) == 0 {
		return &balances, nil
	}
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
	w.Balances = data
	return nil
}
