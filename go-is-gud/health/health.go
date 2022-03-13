package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func New(lastTypedAt *time.Time, port string) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Status(200).JSON(&fiber.Map{
			"lastTypedAt": lastTypedAt.Format(time.RFC3339),
			"status":      "Ok",
		})
	})
	app.Listen(":" + port)
}
