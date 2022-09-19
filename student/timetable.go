package student

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    periods, err := grip.GetPeriods(
      base.Query {
        {"subject.grade.key": student.Grade.Key},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })
}