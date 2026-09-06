package routes

import (
	fiber "github.com/gofiber/fiber/v3"

	"soalhub/internal/modules/user"
)

func AuthRoutes(router fiber.Router, userHandler *user.UserHandler) {
	auth := router.Group("/auth")

	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)
}
