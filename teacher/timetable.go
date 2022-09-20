package teacher

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    var periods []grip.Period
    for _, subject := range teacher.Subjects {
      subjectPeriods, err := grip.GetPeriods(
        base.Query {
          {"subject.key": subject.Key},
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }

      periods = append(periods, subjectPeriods...)
    }

    return c.JSON(periods)
  })
}