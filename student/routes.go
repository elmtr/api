package student

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
  g := app.Group("/v1/student")

  signup(g)
  login(g)
  get(g)
  timetable(g)
}