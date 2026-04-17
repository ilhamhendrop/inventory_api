package util

import (
	"fmt"
	"inventory-app/internal/dto"

	"github.com/gofiber/fiber/v2"
)

func InternalServiceError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ResponseError(err.Error()))
}

func UnprocessableEntity(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dto.ResponseError(err.Error()))
}

func OK(ctx *fiber.Ctx, data any) error {
	return ctx.Status(fiber.StatusOK).JSON(dto.ResponseSucces(data))
}

func Created(ctx *fiber.Ctx, data any) error {
	return ctx.Status(fiber.StatusCreated).JSON(dto.ResponseCreated(data))
}

func BadRequest(ctx *fiber.Ctx, err any) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(dto.ResponseError(fmt.Sprint(err)))
}

func Unauthorized(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(dto.ResponseError(err.Error()))
}

func NoContent(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNoContent).JSON(dto.ResponseNoContent())
}
