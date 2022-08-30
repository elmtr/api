package admin

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
  g := app.Group("/v1/admin")
  
  g.Get("/test", func (c *fiber.Ctx) error {
    return c.SendString("how are you doing?")
  })

  login(g)
  student(g)
  teacher(g)
}