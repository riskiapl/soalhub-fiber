package middleware

import (
	"strings"

	fiber "github.com/gofiber/fiber/v3"

	"soalhub/pkg/utils"
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid Authorization header format")
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		if claims.TokenTYpe != "access" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token type")
		}

		c.Locals("userId", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied. Admins only.")
		}

		return c.Next()
	}
}
