package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	authService model.AuthService
}

func NewAuth(app *fiber.App, authService model.AuthService) {
	aa := authApi{
		authService: authService,
	}

	app.Post("/login", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	}, aa.Login)
}

func (aa authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return util.BadRequest(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	res, err := aa.authService.Login(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}
