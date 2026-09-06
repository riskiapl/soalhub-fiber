package routes

import (
	fiber "github.com/gofiber/fiber/v3"

	"soalhub/internal/middleware"
	"soalhub/internal/modules/user"
)

func AdminRoutes(router fiber.Router, userHandler *user.UserHandler) {
	admin := router.Group("/admin", middleware.RolesAllowed("admin"))

	adminUser := admin.Group("/users")
	adminUser.Get("", userHandler.GetAllUsers)
	adminUser.Get("/:id", userHandler.GetUserByID)
	adminUser.Put("/:id", userHandler.UpdateUser)
	adminUser.Delete("/:id", userHandler.DeleteUser)
}
