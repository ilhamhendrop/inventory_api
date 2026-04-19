package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type financeApi struct {
	financeService model.FinanceService
}

func NewFinance(app *fiber.App, financeService model.FinanceService, authMid fiber.Handler) {
	fa := financeApi{
		financeService: financeService,
	}

	finance := app.Group("/finances", authMid, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	finance.Get("", util.AdminOnly, fa.Index)
	finance.Get("/search", util.AdminOnly, fa.Search)
	finance.Get("/:id", util.AdminOnly, fa.Detail)
	app.Post("/maintenances/:maintenanceId/finances", authMid, util.AdminOnly, fa.Create)
	finance.Patch("/:id", util.AdminOnly, fa.Update)
	finance.Delete("/:id", util.AdminOnly, fa.Delete)
}

func (fa financeApi) Index(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	maintenanceId := ctx.Query("maintenace_id", "")

	res, err := fa.financeService.Index(c, model.FinanceSearch{
		MaintenanceId: maintenanceId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (fa financeApi) Search(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	maintenanceId := ctx.Query("maintenance_id", "")
	keyword := ctx.Query("q")
	if len(keyword) == 0 {
		return util.BadRequest(ctx, keyword)
	}

	res, err := fa.financeService.Search(c, keyword, model.FinanceSearch{
		MaintenanceId: maintenanceId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (fa financeApi) Detail(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	maintenanceId := ctx.Query("maintenance_id", "")
	id := ctx.Params("id")

	res, err := fa.financeService.Detail(c, id, model.FinanceSearch{
		MaintenanceId: maintenanceId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (fa financeApi) Create(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.FinanceCreated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	userId := util.GetUserId(ctx)
	maintenanceId := ctx.Params("maintenance_id")

	err := fa.financeService.Create(c, userId, maintenanceId, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.Created(ctx, "")
}

func (fa financeApi) Update(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.FinanceUpdated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := fa.financeService.Update(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (fa financeApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")

	err := fa.financeService.Delete(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}
