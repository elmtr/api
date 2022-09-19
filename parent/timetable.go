package parent

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var parent grip.Parent
    utils.GetLocals(c, "parent", &parent)

    student, err := grip.GetStudent(
      base.Query {
        {"key": c.Query("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

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