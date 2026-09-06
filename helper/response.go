package helper

import (
	"api-students/app/model"
	"github.com/gofiber/fiber/v2"
)

func Ok(c *fiber.Ctx, msg string, data any) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

func OkList(c *fiber.Ctx, msg string, data any, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *fiber.Ctx, msg string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Fail(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(model.WebResponse{
		Status:  "fail",
		Message: msg,
	})
}

func FailValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{
		Status:  "fail",
		Message: "validasi gagal",
		Errors:  errs,
	})
}