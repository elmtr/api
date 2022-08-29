package teacher

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/v1/teacher")

  signup(g)
  login(g)
  get(g)
  set(g)
  update(g)
}