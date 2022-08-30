package admin

import (
	"api/bong"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func student(g fiber.Router) {
  student := g.Group("/student")

  student.Get("", authMiddleware, func (c *fiber.Ctx) error {
    phone := c.Query("phone")

    student, err := bong.GetStudent(
      bson.M{
        "phone": phone,
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    
    return c.JSON(student)
  })

  student.Post("/subjects/add", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    student, err := bong.GetStudent(
      bson.M{
        "id": body["id"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := bong.GetSubject(
      bson.M{
        "grade.id": student.Grade.ID,
        "name": body["name"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    shortSubject := bong.ShortSubject {
      ID: subject.ID,
      Name: subject.Name,
    }

    student.Subjects, err = bong.AddStudentSubject(student.ID, student.Subjects, shortSubject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(student)
  })

  student.Post("/subjects/remove", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    student, err := bong.GetStudent(
      bson.M{
        "id": body["id"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := bong.GetSubject(
      bson.M{
        "grade.id": student.Grade.ID,
        "name": body["name"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    shortSubject := bong.ShortSubject {
      ID: subject.ID,
      Name: subject.Name,
    }

    student.Subjects, err = bong.RemoveStudentSubject(student.ID, student.Subjects, shortSubject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(student)
  })

  student.Patch("/grade", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

    grade, err := bong.GetGrade(
      bson.M{
        "gradeNumber": gradeNumber,
        "gradeLetter": body["gradeLetter"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    student, err := bong.UpdateStudentGrade(body["id"], grade)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(student)
  })
}