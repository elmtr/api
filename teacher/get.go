package teacher

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func get(g fiber.Router) {
  g.Get("/students", authMiddleware, func (c *fiber.Ctx) error {
    students, err := grip.GetStudents(
      base.Query {
        {"grade.key": c.Query("gradeKey")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(students)
  })

  g.Get("/marks", authMiddleware, func (c *fiber.Ctx) error {
    marks, err := grip.GetMarks(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": c.Query("studentKey"),
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(marks)
  })

  g.Get("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    truancies, err := grip.GetTruancies(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": c.Query("studentKey"),
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(truancies)
  })

  g.Get("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    draftMarks, err := grip.GetDraftMarks(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": c.Query("studentKey"),
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(draftMarks)
  })

  g.Get("/points", authMiddleware, func (c *fiber.Ctx) error {
    points, err := grip.GetPoints(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": c.Query("studentKey"),
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(points)
  })
}