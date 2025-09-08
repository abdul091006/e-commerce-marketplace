package utils

// UpdateBalanceRequest represents the request to update wallet balance
type UpdateBalanceRequest struct {
	BalanceType string `json:"type" validate:"required"`
	Amount      string  `json:"amount" validate:"required,numeric"`
}