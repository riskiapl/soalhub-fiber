package routes

import (
	fiber "github.com/gofiber/fiber/v3"

	"soalhub/internal/middleware"
)

func StudentRoutes(router fiber.Router) {
	student := router.Group("/student", middleware.RolesAllowed("student, admin"))

	_ = student
}
