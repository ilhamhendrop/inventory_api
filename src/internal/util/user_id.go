package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserId(ctx *fiber.Ctx) string {
	user, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || user == nil {
		return ""
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	if v, ok := claims["id"].(string); ok {
		return v
	}

	return ""
}
