package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return fiber.ErrUnauthorized
		}

		token, ok := user.(*jwt.Token)
		if !ok {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		role, ok := claims["role"].(string)
		if !ok {
			return fiber.ErrForbidden
		}

		role = strings.ToLower(role)

		for _, v := range roles {
			if strings.ToLower(v) == role {
				return c.Next()
			}
		}

		return fiber.ErrForbidden
	}
}
