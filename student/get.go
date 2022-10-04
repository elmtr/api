package student

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func get(g fiber.Router) {
  g.Get("/subjects", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    subjects, err := grip.GetSubjects(
      base.Query {
        {
          "grade.gradeKey": student.Grade.Key,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(subjects)
  })

  g.Get("/marks", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    marks, err := grip.GetMarks(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": student.Key,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(marks)
  })

  g.Get("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    truancies, err := grip.GetTruancies(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": student.Key,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(truancies)
  })

  g.Get("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    draftMarks, err := grip.GetDraftMarks(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": student.Key,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(draftMarks)
  })

  g.Get("/points", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    points, err := grip.GetPoints(
      base.Query {
        {
          "subjectKey": c.Query("subjectKey"),
          "studentKey": student.Key,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(points)
  })
}