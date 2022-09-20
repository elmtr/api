package teacher

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func homeroom(g fiber.Router) {
  hr := g.Group("/homeroom")

  hr.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    students, err := grip.GetStudents(
      base.Query {
        {"grade.key": teacher.Homeroom.Key},
      },
    )
    if err != nil {
      utils.Error(c, err)
    }

    return c.JSON(students)
  })

  hr.Get("/truancies", func (c *fiber.Ctx) error {
    key := c.Query("studentKey")
    
    truancies, err := grip.GetTruancies(
      base.Query {
        {"studentKey": key},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(truancies)
  })

  hr.Get("/timetable", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    periods, err := grip.GetPeriods(
      base.Query {
        {"subject.grade.key": teacher.Homeroom.Key},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })

  // TODO: be able to close a year average mark: after the law gets clearer
}
