package utils

import (
	"errors"
	"fmt"
)

// Custom error types
var (
	ErrWalletNotFound        = errors.New("wallet not found")
	ErrWalletAlreadyExists   = errors.New("wallet already exists for this user")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInvalidBalanceType    = errors.New("invalid balance type")
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidWalletUserID   = errors.New("invalid wallet user ID")
	ErrDatabaseOperation     = errors.New("database operation failed")
	ErrJSONUnmarshal         = errors.New("failed to unmarshal JSON")
	ErrJSONMarshal           = errors.New("failed to marshal JSON")
)

// WalletError represents a wallet service specific error
type WalletError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *WalletError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewWalletError creates a new WalletError
func NewWalletError(code, message, details string) *WalletError {
	return &WalletError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Error codes
const (
	CodeWalletNotFound      = "WALLET_NOT_FOUND"
	CodeWalletExists        = "WALLET_ALREADY_EXISTS"
	CodeInsufficientBalance = "INSUFFICIENT_BALANCE"
	CodeInvalidAmount       = "INVALID_AMOUNT"
	CodeInvalidBalanceType  = "INVALID_BALANCE_TYPE"
	CodeDatabaseError       = "DATABASE_ERROR"
	CodeValidationError     = "VALIDATION_ERROR"
	CodeInternalError       = "INTERNAL_ERROR"
)

// IsWalletError checks if an error is a WalletError
func IsWalletError(err error) bool {
	_, ok := err.(*WalletError)
	return ok
}

// GetErrorCode returns the error code from a WalletError
func GetErrorCode(err error) string {
	if walletErr, ok := err.(*WalletError); ok {
		return walletErr.Code
	}
	return CodeInternalError
}