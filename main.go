package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// Middleware untuk menolak request yang bukan JSON
func requireJSON(c *fiber.Ctx) error {
	if metodeBerbody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	app := fiber.New()

	api := app.Group("/api/v1")
	s := api.Group("/students", requireJSON)

	// Routing endpoint POST
	s.Post("/", createStudent)
	// Routing GET
	s.Get("/", listStudents)
	s.Get("/:id", getStudent)
	s.Put("/:id", replaceStudent)
	s.Patch("/:id", patchStudent)
	s.Delete("/:id", deleteStudent)

	log.Println("Server jalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
