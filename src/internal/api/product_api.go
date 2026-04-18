package api

import (
	"inventory-app/internal/dto"
	"inventory-app/internal/model"
	"inventory-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type productApi struct {
	productService model.ProductService
}

func NewProduct(app *fiber.App, productService model.ProductService, authMid fiber.Handler) {
	pa := productApi{
		productService: productService,
	}

	product := app.Group("/products", authMid, func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.Next()
	})

	product.Get("", util.AdminOnly, pa.Index)
	product.Get("/search", util.AdminOnly, pa.Search)
	product.Get("/:id", util.AdminOnly, pa.Detail)
	product.Post("", util.AdminOnly, pa.Create)
	product.Patch("/:id", util.AdminOnly, pa.Update)
	product.Delete("/:id", util.AdminOnly, pa.Delete)
}

func (pa productApi) Index(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	categorieId := ctx.Query("categorie_id", "")

	res, err := pa.productService.Index(c, model.ProductSearch{
		CategorieId: categorieId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (pa productApi) Search(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	categorieId := ctx.Query("categorie_id", "")
	keyword := ctx.Query("q")
	if len(keyword) == 0 {
		return util.BadRequest(ctx, keyword)
	}

	res, err := pa.productService.Search(c, keyword, model.ProductSearch{
		CategorieId: categorieId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (pa productApi) Detail(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	categorieId := ctx.Query("categorie_id", "")
	id := ctx.Params("id")

	res, err := pa.productService.Detail(c, id, model.ProductSearch{
		CategorieId: categorieId,
	})
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, res)
}

func (pa productApi) Create(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.ProductCreated
	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	err := pa.productService.Create(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.Created(ctx, "")
}

func (pa productApi) Update(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	var req dto.ProductUpdated

	if err := ctx.BodyParser(&req); err != nil {
		return util.UnprocessableEntity(ctx, err)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return util.BadRequest(ctx, fails)
	}

	req.ID = ctx.Params("id")
	err := pa.productService.Update(c, req)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}

func (pa productApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := util.WithTimeout(ctx)
	defer cancel()

	id := ctx.Params("id")
	err := pa.productService.Delete(c, id)
	if err != nil {
		return util.InternalServiceError(ctx, err)
	}

	return util.OK(ctx, "")
}
