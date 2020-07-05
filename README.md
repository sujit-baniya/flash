# Set flash message for routes.

This package is build to send the flash messages on the top of Gofiber

## Installation
The package can be used to validate the data and send flash message to other route.
> go get github.com/itsursujit/flash


## Usage

```go
package main

import (
    "github.com/gofiber/fiber"
    "github.com/itsursujit/flash"
)

func main() {
    app := fiber.New()

    f := flash.Flash{
        CookiePrefix:"Test-Prefix",
    }
    app.Get("/success-redirect", func (c *fiber.Ctx) {
        f.Get(c)
        c.JSON(f.Data)
    })

    app.Get("/success", func (c *fiber.Ctx) {
        mp := fiber.Map{
            "success": true,
            "message": "I'm receiving success",
        }
        c.Redirect("/success-redirect")
    })

    app.Get("/error-redirect", func (c *fiber.Ctx) {
        f.Get(c)
        c.JSON(f.Data)
    })

    app.Get("/error", func (c *fiber.Ctx) {
        mp := fiber.Map{
            "error": true,
            "message": "I'm receiving error",
        }
        f.Data = mp
        f.Error(c)
        c.Redirect("/error-redirect")
    })

    app.Listen(3000)
}
```