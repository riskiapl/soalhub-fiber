package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"soalhub/internal/modules/user"
)

func Connect() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pg_trgm;").Error; err != nil {
		log.Fatalf("Gagal mengaktifkan ekstensi pg_trgm: %v", err)
	}

	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_users_name_trgm ON users USING gin (name gin_trgm_ops);`
	if err := db.Exec(createIndexSQL).Error; err != nil {
		log.Printf("Warning: Failed to create GIN trigram index on users.name: %v", err)
	}

	SeedAdmin(db)

	log.Println("Connected to the database successfully")
	return db
}
