package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"api-students/app/model"
)

// reqCtx memberi batas waktu 5 detik untuk setiap operasi basis data
func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func ok(c *fiber.Ctx, msg string, data any) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

func okList(c *fiber.Ctx, msg string, data any, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

func created(c *fiber.Ctx, msg string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(model.WebResponse{
		Status:  "fail",
		Message: msg,
	})
}

func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{
		Status:  "fail",
		Message: "validasi gagal",
		Errors:  errs,
	})
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	return id, err == nil && id > 0
}

func parseListQuery(c *fiber.Ctx) model.ListQuery {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	sort := c.Query("sort", "id")
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		val := activeStr == "true"
		isActive = &val
	}

	return model.ListQuery{
		Page:     page,
		Limit:    limit,
		Search:   search,
		Sort:     sort,
		Order:    order,
		IsActive: isActive,
	}
}