package routes

import (
	fiber "github.com/gofiber/fiber/v3"

	"soalhub/internal/middleware"
)

func TeacherRoutes(router fiber.Router) {
	teacher := router.Group("/teacher", middleware.RolesAllowed("teacher, admin"))

	_ = teacher
}
