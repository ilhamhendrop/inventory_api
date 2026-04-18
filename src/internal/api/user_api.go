package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type userApi struct {
	userService model.UserService
}

func NewUser(app *fiber.App, userService model.UserService, authMid fiber.Handler) {
	ua := userApi{
		userService: userService,
	}

	user := app.Group("/users", authMid, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	user.Get("", util.AdminOnly, ua.Index)
	user.Get("/search", util.AdminOnly, ua.Search)
	user.Get("/:id", util.AdminOnly, ua.Detail)
	user.Post("", authMid, ua.Create)
	user.Patch("/:id/data", util.AdminOnly, ua.UpdateData)
	user.Patch("/:id/password", util.AdminOnly, ua.UpdatePassword)
	user.Delete("/:id", util.AdminOnly, ua.Delete)
}

func (ua userApi) Index(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	res, err := ua.userService.Index(c)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ua userApi) Search(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	keyword := ctx.Query("q")
	if len(keyword) == 0 {
		return util.BadRequest(ctx, keyword)
	}

	res, err := ua.userService.Search(c, keyword)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ua userApi) Detail(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")
	res, err := ua.userService.Detail(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ua userApi) Create(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.UserCreated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	err := ua.userService.Create(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.Created(ctx, "")
}

func (ua userApi) UpdateData(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.UserUpdateData
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := ua.userService.UpdateData(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (ua userApi) UpdatePassword(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.UserUpdatePassword
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := ua.userService.UpdatePassword(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (ua userApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")
	err := ua.userService.Delete(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}
