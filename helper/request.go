package helper

import (
	"context"
	"strconv"
	"strings"
	"time"

	"api-students/app/model"
	"github.com/gofiber/fiber/v2"
)

// ReqCtx memberi batas waktu 5 detik untuk setiap operasi basis data
func ReqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func ParamID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	return id, err == nil && id > 0
}

func ParseListQuery(c *fiber.Ctx) model.ListQuery {
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