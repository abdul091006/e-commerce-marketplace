package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// SuccessResponse returns a successful API response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// CreatedResponse returns a created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
	})
}

// BadRequestResponse returns a bad request response
func BadRequestResponse(c *fiber.Ctx, message string, err interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, err)
}

// NotFoundResponse returns a not found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

// InternalServerErrorResponse returns an internal server error response
func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, nil)
}

// ConflictResponse returns a conflict response
func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, message, nil)
}