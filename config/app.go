package config

import (
	"log/slog"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
	"api-students/route"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApp(logger *slog.Logger, pool *pgxpool.Pool, jwt *helper.JWTManager, studentService *service.StudentService, prestasiService *service.PrestasiService, authService *service.AuthService) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "Praktikum Backend Lanjut"),
		ErrorHandler: newErrorHandler(logger),
	})

	middleware.Register(app, logger)
	route.Register(app, route.Dependencies{
		Pool:            pool,
		JWT:             jwt,
		StudentService:  studentService,
		PrestasiService: prestasiService,
		AuthService:     authService,
	})
	
	app.Use(func(c *fiber.Ctx) error {
		return helper.Fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		message := "terjadi error pada server"

		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
			message = e.Message
		}

		logger.Error("unhandled_error",
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.String("error", err.Error()),
		)

		return helper.Fail(c, status, message)
	}
}