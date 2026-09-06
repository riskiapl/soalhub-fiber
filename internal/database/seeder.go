package database

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"soalhub/internal/modules/user"
)

// SeedAdmin membuat akun admin default jika belum ada di database
func SeedAdmin(db *gorm.DB) {
	adminEmail := "admin@soalhub.com"

	// Cek apakah admin sudah ada
	var count int64
	db.Model(&user.User{}).Where("email = ?", adminEmail).Count(&count)

	if count > 0 {
		return // Admin sudah ada, tidak perlu re-seed
	}

	// Hash password admin
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "default_password"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	// Buat entitas admin baru
	adminUser := user.User{
		Name:     "Admin",
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Printf("⚠️ Failed to seed admin account: %v\n", err)
		return
	}

	log.Println("🔑 Admin account created successfully: admin@soalhub.com / xml5nfGl3F0n713FqReK")
}
