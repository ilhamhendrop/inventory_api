package util

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const defaultTimeout = 10 * time.Second

func WithTimeout(ctx *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Context(), defaultTimeout)
}
