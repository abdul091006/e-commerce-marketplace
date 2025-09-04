package utils

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError

	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			switch err.Tag() {
			case "required":
				message = "This field is required"
			case "min":
				message = "Value is too small"
			case "max":
				message = "Value is too large"
			case "oneof":
				message = "Invalid value. Allowed values: " + err.Param()
			default:
				message = "Invalid value"
			}

			validationErrors = append(validationErrors, ValidationError{
				Field:   ToSnakeCase(err.Field()),
				Message: message,
			})
		}
	}

	return validationErrors
}

// ValidateWalletUserID validates wallet user ID format
func ValidateWalletUserID(walletUserID string) error {
	if walletUserID == "" {
		return NewWalletError(CodeValidationError, "wallet_user_id is required", "")
	}

	if len(walletUserID) > 255 {
		return NewWalletError(CodeValidationError, "wallet_user_id is too long", "maximum 255 characters allowed")
	}

	// Basic validation - alphanumeric, hyphens, underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, walletUserID)
	if !matched {
		return NewWalletError(CodeValidationError, "invalid wallet_user_id format", "only alphanumeric characters, hyphens, and underscores are allowed")
	}

	return nil
}

// ValidateBalanceType validates balance type
func ValidateBalanceType(balanceType string) error {
	validTypes := map[string]bool{
		"coins": true,
		"exp":   true,
	}

	if !validTypes[balanceType] {
		return NewWalletError(CodeInvalidBalanceType, "invalid balance type", "allowed types: coins, exp")
	}

	return nil
}

// ValidateAmount validates amount value (float64)
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return NewWalletError(CodeInvalidAmount, "amount must be positive", "")
	}

	if amount > 1000000000 { // 1 billion limit
		return NewWalletError(CodeInvalidAmount, "amount is too large", "maximum 1,000,000,000 allowed")
	}

	return nil
}


// ToSnakeCase converts camelCase to snake_case
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
