package admin

import (
	"api/grip"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func student(g fiber.Router) {
  student := g.Group("/student")

  student.Get("", authMiddleware, func (c *fiber.Ctx) error {
    phone := c.Query("phone")
    key := c.Query("key")

    // getting one student by 
    if (phone == "") {
      student, err := grip.GetStudent(
        base.Query {
          {"key": key},
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
      return c.JSON(student)
    } else {
      student, err := grip.GetStudent(
        base.Query {
          {"phone": phone},
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
      return c.JSON(student)
    }
  })

  student.Get("/grade", authMiddleware, func (c *fiber.Ctx) error {
    gradeKey := c.Query("gradeKey")

    // getting students by gradeKey 
    students, err := grip.GetStudents(
      base.Query {
        {"grade.key": gradeKey},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(students)
  })

  // student.Post("/subjects/add", authMiddleware, func (c *fiber.Ctx) error {
  //   var body map[string]string
  //   json.Unmarshal(c.Body(), &body)

  //   student, err := grip.GetStudent(
  //     base.Query {
  //       {"key": body["key"]},
  //     },
  //   )
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }

  //   subject, err := grip.GetSubject(
  //     base.Query{
  //       {
  //         "grade.key": student.Grade.Key,
  //         "name": body["name"],
  //       },
  //     },
  //   )
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }
  //   shortSubject := grip.ShortSubject {
  //     ID: subject.ID,
  //     Name: subject.Name,
  //   }

  //   student.Subjects, err = grip.AddStudentSubject(student.ID, student.Subjects, shortSubject)
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }

  //   return c.JSON(student)
  // })

  // student.Post("/subjects/remove", authMiddleware, func (c *fiber.Ctx) error {
  //   var body map[string]string
  //   json.Unmarshal(c.Body(), &body)

  //   student, err := grip.GetStudent(
  //     base.Query {
  //       "id": body["id"],
  //     },
  //   )
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }

  //   subject, err := grip.GetSubject(
  //     base.Query {
  //       "grade.id": student.Grade.ID,
  //       "name": body["name"],
  //     },
  //   )
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }
  //   shortSubject := grip.ShortSubject {
  //     ID: subject.ID,
  //     Name: subject.Name,
  //   }

  //   student.Subjects, err = grip.RemoveStudentSubject(student.ID, student.Subjects, shortSubject)
  //   if err != nil {
  //     return utils.Error(c, err)
  //   }

  //   return c.JSON(student)
  // })

  student.Patch("/grade", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

    grade, err := grip.GetGrade(
      base.Query {
        {
          "gradeNumber": gradeNumber,
          "gradeLetter": body["gradeLetter"],
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    err = grip.UpdateStudentGrade(body["key"], grade)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(grade)
  })
}