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
	AuthRoutes(api, userHandler)

	// =============== User Management Routes ===============
	protected := api.Group("/users", middleware.Protected())

	// =============== Admin Routes ===============
	AdminRoutes(protected, userHandler)

	// =============== Teacher Routes ===============
	TeacherRoutes(protected)

	// =============== Student Routes ===============
	StudentRoutes(protected)
}
