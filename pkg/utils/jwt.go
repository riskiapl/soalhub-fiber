package utils

import (
	"errors"
	"os"
	"soalhub/internal/modules/user"
	"time"

	fiber "github.com/gofiber/fiber/v3"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID    uint      `json:"user_id"`
	Role      user.Role `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "soalhub_secret_key_change_me_in_production"
	}
	return []byte(secret)
}

func GenerateAccessToken(userID uint, role user.Role) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func GenerateRefreshToken(userID uint, role user.Role) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetCurrentUserID(c fiber.Ctx) (uint, user.Role) {
	userID, _ := c.Locals("userID").(uint)
	role, _ := c.Locals("role").(user.Role)

	return userID, role
}
