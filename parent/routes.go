package parent

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/v1/parent")

  signup(g)
  login(g)
  get(g)
  update(g)
}