package user

import (
	fiber "github.com/gofiber/fiber/v3"

	"soalhub/pkg/utils"
)

type UserHandler struct {
	service UserService
}

func NewHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// ================ Authentication Handlers ================
func (h *UserHandler) Login(c fiber.Ctx) error {
	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, utils.FormatValidationError(err))
	}

	res, err := h.service.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", res)
}

func (h *UserHandler) RefreshToken(c fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, utils.FormatValidationError(err))
	}

	res, err := h.service.RefreshToken(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", res)
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, utils.FormatValidationError(err))
	}

	res, err := h.service.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", res)
}

// ================ User Management Handlers ================
func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	userID := c.Params("id")

	return utils.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", fiber.Map{
		"id":   userID,
		"name": "John Doe",
	})
}

func (h *UserHandler) GetAllUsers(c fiber.Ctx) error {
	var param UserQueryParam

	if err := c.Bind().Query(&param); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.GetAllUsers(param)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "All users retrieved successfully", res)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
