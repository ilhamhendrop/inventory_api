package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type warehouseApi struct {
	warehouseService model.WarehouseService
}

func NewWarehouse(app *fiber.App, warehouseService model.WarehouseService, authMid fiber.Handler) {
	wa := warehouseApi{
		warehouseService: warehouseService,
	}

	warehouse := app.Group("/warehouses", authMid, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	warehouse.Get("", util.AdminOnly, wa.Index)
	warehouse.Get("/search", util.AdminOnly, wa.Search)
	warehouse.Get("/:id", util.AdminOnly, wa.Detail)
	warehouse.Post("", util.AdminOnly, wa.Create)
	warehouse.Patch("/:id", util.AdminOnly, wa.Update)
	warehouse.Delete("/:id", util.AdminOnly, wa.Delete)
}

func (wa warehouseApi) Index(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	productId := ctx.Query("product_id", "")
	res, err := wa.warehouseService.Index(c, model.WarehouseSearch{
		ProductId: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (wa warehouseApi) Search(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	productId := ctx.Query("product_id", "")
	keyword := ctx.Query("q")
	if len(keyword) == 0 {
		return util.BadRequest(ctx, keyword)
	}

	res, err := wa.warehouseService.Search(c, keyword, model.WarehouseSearch{
		ProductId: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (wa warehouseApi) Detail(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	productId := ctx.Query("product_id", "")
	id := ctx.Params("id")

	res, err := wa.warehouseService.Detail(c, id, model.WarehouseSearch{
		ProductId: productId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (wa warehouseApi) Create(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.WarehouseCreated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	err := wa.warehouseService.Create(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.Created(ctx, "")
}

func (wa warehouseApi) Update(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.WarehouseUpdated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := wa.warehouseService.Update(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (wa warehouseApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")
	err := wa.warehouseService.Delete(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}
