package utils

import (
	"errors"
	"strconv"

	fiber "github.com/gofiber/fiber/v3"
)

type AuthUser struct {
	UserID uint
	Role   string
}

// ExtractAuthUser mengambil ID dan Role pengguna yang sedang login dari JWT context
func ExtractAuthUser(c fiber.Ctx) (AuthUser, error) {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return AuthUser{}, errors.New("unauthorized: missing or invalid user ID")
	}

	role, ok := c.Locals("role").(string)
	if !ok {
		return AuthUser{}, errors.New("unauthorized: missing or invalid role")
	}

	return AuthUser{
		UserID: userID,
		Role:   role,
	}, nil
}

// ExtractTargetIDAndAuth mengambil param ID target dari URL (misal: /:id) sekaligus data Auth user
func ExtractTargetIDAndAuth(c fiber.Ctx, paramName string) (uint, AuthUser, error) {
	idParam := c.Params(paramName)
	targetID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, AuthUser{}, errors.New("invalid target ID format")
	}

	authUser, err := ExtractAuthUser(c)
	if err != nil {
		return 0, AuthUser{}, err
	}

	return uint(targetID), authUser, nil
}
