# Set flash message for routes.

This package is build to send the flash messages on the top of Gofiber

## Installation
The package can be used to validate the data and send flash message to other route.
> go get github.com/sujit-baniya/flash


## Usage

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/sujit-baniya/flash"
)

func main() {
	app := fiber.New()

	f := flash.Flash{
		CookiePrefix:"Test-Prefix",
	}
	app.Get("/success-redirect", func (c *fiber.Ctx) error {
		f.Get(c)
		return c.JSON(f.Data)
	})

	app.Get("/success", func (c *fiber.Ctx) error {
		mp := fiber.Map{
			"success": true,
			"message": "I'm receiving success",
		}
		f.Data = mp
		f.Success(c)
		return c.Redirect("/success-redirect")
	})

	app.Get("/error-redirect", func (c *fiber.Ctx) error {
		f.Get(c)
		return c.JSON(f.Data)
	})

	app.Get("/error", func (c *fiber.Ctx) error {
		mp := fiber.Map{
			"error": true,
			"message": "I'm receiving error",
		}
		f.Data = mp
		f.Error(c)
		return c.Redirect("/error-redirect")
	})

	app.Get("/error-with-data", func (c *fiber.Ctx) error {
		mp := fiber.Map{
			"error": true,
			"message": "I'm receiving error with inline error data",
		}
		return f.WithError(c, mp).Redirect("/error-redirect")
	})

	app.Get("/success-with-data", func (c *fiber.Ctx) error {
		mp := fiber.Map{
			"success": true,
			"message": "I'm receiving success with inline success data",
		}
		return f.WithSuccess(c, mp).Redirect("/success-redirect")
	})

	app.Listen(3000)
}
```
