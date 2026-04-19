package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type maintenanceApi struct {
	maintenanceService model.MaintenanceService
}

func NewMaintenance(app *fiber.App, maintenanceService model.MaintenanceService, authMid fiber.Handler) {
	ma := maintenanceApi{
		maintenanceService: maintenanceService,
	}

	maintenance := app.Group("/maintenances", authMid, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	maintenance.Get("", util.AdminOnly, ma.Index)
	maintenance.Get("/search", util.AdminOnly, ma.Search)
	maintenance.Get("/:id", util.AdminOnly, ma.Detail)
	app.Post("/warehouses/:warehouseId/maintenances", authMid, util.AdminOnly, ma.Create)
	maintenance.Patch("/:id", util.AdminOnly, ma.Update)
	maintenance.Delete("/:id", util.AdminOnly, ma.Delete)
}

func (ma maintenanceApi) Index(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	productId := ctx.Query("product_id", "")

	res, err := ma.maintenanceService.Index(c, model.MaintenanceSearch{
		ProductID: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ma maintenanceApi) Search(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	productId := ctx.Query("product_id", "")
	keyword := ctx.Query("q")
	if len(keyword) == 0 {
		return util.BadRequest(ctx, keyword)
	}

	res, err := ma.maintenanceService.Search(c, keyword, model.MaintenanceSearch{
		ProductID: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ma maintenanceApi) Detail(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")
	productId := ctx.Query("product_id", "")

	res, err := ma.maintenanceService.Detail(c, id, model.MaintenanceSearch{
		ProductID: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (ma maintenanceApi) Create(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.MaintenanceCreated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	userId := util.GetUserId(ctx)
	warehouseId := ctx.Params("warehouse_id")

	err := ma.maintenanceService.Create(c, userId, warehouseId, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.Created(ctx, "")
}

func (ma maintenanceApi) Update(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.MaintenanceUpdated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := ma.maintenanceService.Update(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (ma maintenanceApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")

	err := ma.maintenanceService.Delete(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}
