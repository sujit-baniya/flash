package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

func main() {
	app := fiber.New()
	app.Get("/success-redirect", func(c *fiber.Ctx) error {
		return c.JSON(flash.Get(c))
	})

	app.Get("/error-redirect", func(c *fiber.Ctx) error {
		flash.Get(c)
		return c.JSON(flash.Get(c))
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		mp := fiber.Map{
			"error":   true,
			"message": "I'm receiving error with inline error data",
		}
		return flash.WithError(c, mp).Redirect("/error-redirect")
	})

	app.Get("/success", func(c *fiber.Ctx) error {
		mp := fiber.Map{
			"success": true,
			"message": "I'm receiving success with inline success data",
		}
		return flash.WithSuccess(c, mp).Redirect("/success-redirect")
	})

	app.Get("/data", func(c *fiber.Ctx) error {
		mp := fiber.Map{
			"text": "Received arbitrary data",
		}
		return flash.WithData(c, mp).Redirect("/success-redirect")
	})

	app.Listen(":8080")
}
