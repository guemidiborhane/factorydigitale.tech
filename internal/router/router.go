package router

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(a *fiber.App) {
	a.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Howdy 👋",
		})
	})
	setupAPIRoutes(a)
}
