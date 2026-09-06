package main

import (
	"log"
	"os"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	"soalhub/internal/database"
	"soalhub/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db := database.Connect()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	app := fiber.New(fiber.Config{
		AppName: "SoalHub API",
	})

	// Health check endpoint
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "SoalHub API is running",
		})
	})

	routes.SetupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified in .env
	}
	log.Fatal(app.Listen(":"+port, fiber.ListenConfig{
		// enable on production for better performance on multi-core systems
		// EnablePrefork: true,
	}))
}
