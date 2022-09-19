package admin

import (
	"api/grip"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func teacher(g fiber.Router) {
  teacher := g.Group("/teacher")

  teacher.Get("", authMiddleware, func (c *fiber.Ctx) error {
    phone := c.Query("phone")

    teacher, err := grip.GetTeacher(
      base.Query {
        {"phone": phone},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    
    return c.JSON(teacher)
  })

  teacher.Post("/subjects/add", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

    teacher, err := grip.GetTeacher(
      base.Query {
        {"key": body["key"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := grip.GetSubject(
      base.Query {
        {
          "name": body["name"],
          "grade.gradeNumber": gradeNumber,
          "grade.gradeLetter": body["gradeLetter"],
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    teacher.Subjects, err = grip.AddTeacherSubject(teacher.Key, teacher.Subjects, subject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })

  teacher.Post("/subjects/remove", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

    teacher, err := grip.GetTeacher(
      base.Query {
        {"key": body["key"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := grip.GetSubject(
      base.Query {
        {
          "name": body["name"],
          "grade.gradeNumber": gradeNumber,
          "grade.gradeLetter": body["gradeLetter"],
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    teacher.Subjects, err = grip.RemoveTeacherSubject(teacher.Key, teacher.Subjects, subject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })

  teacher.Patch("/homeroom", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

    grade, err := grip.GetGrade(
      base.Query{
        {
          "gradeNumber": gradeNumber,
          "gradeLetter": body["gradeLetter"],
        },
      },
    )
    if err != nil {
      return utils.MessageError(c, "nope")
    }

    teacher, err := grip.UpdateTeacherHomeroom(body["key"], grade)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })
}