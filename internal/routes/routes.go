package routes

import (
	fiber "github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"soalhub/internal/middleware"
	"soalhub/internal/modules/user"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)

	userGroup := api.Group("/users", middleware.Protected())

	userGroup.Get("/:id", userHandler.GetUserByID)

	adminGroup := api.Group("/admin", middleware.Protected(), middleware.AdminOnly())

	adminUserGroup := adminGroup.Group("/users")
	adminUserGroup.Get("", userHandler.GetAllUsers)
}
