package admin

import (
	"api/grip"
	"api/utils"
	"encoding/json"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func admin(g fiber.Router) {
  g.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    return c.SendString("hello world")
  })

  g.Get("/grades", authMiddleware, func (c *fiber.Ctx) error {
    var admin grip.Admin
    utils.GetLocals(c, "admin", &admin)
    grades, err := grip.GetGrades(
      base.Query {
        {"schoolKey": admin.SchoolKey},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(grades)
  })

  g.Post("/grades", authMiddleware, func (c *fiber.Ctx) error {
    var grade grip.Grade
    json.Unmarshal(c.Body(), &grade)

    var admin grip.Admin
    utils.GetLocals(c, "admin", &admin)

    grade.SchoolKey = admin.SchoolKey

    grade.Put()

    return c.JSON(grade)
  })

  g.Get("/subjects", authMiddleware, func (c *fiber.Ctx) error {
    subjects, err := grip.GetSubjects(
      base.Query {
        {"grade.key": c.Query("gradeKey")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(subjects)
  })

  g.Post("/subjects", authMiddleware, func (c *fiber.Ctx) error {
    grade, err := grip.GetGrade(
      base.Query {
        {"key": c.Query("gradeKey")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    var subject grip.Subject
    json.Unmarshal(c.Body(), &subject)
    subject.Grade = grade

    subject.Put()

    return c.JSON(subject)
  })
}