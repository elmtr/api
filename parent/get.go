package parent

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func get(g fiber.Router) {
  g.Get("/student", authMiddleware, func (c *fiber.Ctx) error {
    student, err := grip.GetStudent(
      base.Query {
        {"key": c.Query("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(student)
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