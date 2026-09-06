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

	// =============== Authentication Routes ===============
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
	auth.Post("/refresh", userHandler.RefreshToken)

	// =============== User Management Routes ===============
	protected := api.Group("/users", middleware.Protected())

	// =============== Admin Routes ===============
	adminGroup := protected.Group("/admin", middleware.RolesAllowed("admin"))

	adminUserGroup := adminGroup.Group("/users")
	adminUserGroup.Get("", userHandler.GetAllUsers)

	// =============== Teacher Routes ===============
	// teacherGroup := protected.Group("/teacher", middleware.RolesAllowed("teacher"))

	// =============== Student Routes ===============
	// studentGroup := protected.Group("/student", middleware.RolesAllowed("student"))
}
