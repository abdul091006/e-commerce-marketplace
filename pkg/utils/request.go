package utils

// UpdateBalanceRequest represents the request to update wallet balance
type UpdateBalanceRequest struct {
	BalanceType string `json:"balance_type" validate:"required,oneof=coins exp"`
	Amount      int64  `json:"amount" validate:"required,min=1"`
}