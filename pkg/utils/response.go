package utils

import (
	fiber "github.com/gofiber/fiber/v3"
)

// ResponseFormat merepresentasikan struktur standar JSON API
type ResponseFormat struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SuccessResponse mengembalikan response JSON sukses (misal 200, 201)
func SuccessResponse(c fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(ResponseFormat{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengembalikan response JSON gagal (misal 400, 404, 500)
func ErrorResponse(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(ResponseFormat{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}
