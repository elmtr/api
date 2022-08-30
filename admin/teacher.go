package admin

import (
	"api/bong"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func teacher(g fiber.Router) {
  teacher := g.Group("/teacher")

  teacher.Get("", authMiddleware, func (c *fiber.Ctx) error {
    phone := c.Query("phone")

    teacher, err := bong.GetTeacher(
      bson.M{
        "phone": phone,
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

    teacher, err := bong.GetTeacher(
      bson.M{
        "id": body["id"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := bong.GetSubject(
      bson.M{
        "name": body["name"],
        "grade.gradeNumber": gradeNumber,
        "grade.gradeLetter": body["gradeLetter"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    teacher.Subjects, err = bong.AddTeacherSubject(teacher.ID, teacher.Subjects, subject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })

  teacher.Post("/subjects/remove", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])


    teacher, err := bong.GetTeacher(
      bson.M{
        "id": body["id"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    subject, err := bong.GetSubject(
      bson.M{
        "name": body["name"],
        "grade.gradeNumber": gradeNumber,
        "grade.gradeLetter": body["gradeLetter"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    teacher.Subjects, err = bong.RemoveTeacherSubject(teacher.ID, teacher.Subjects, subject)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })

  teacher.Patch("/homeroom", authMiddleware, func (c *fiber.Ctx) error {
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
      return utils.MessageError(c, "nope")
    }

    teacher, err := bong.UpdateTeacherHomeroom(body["id"], grade)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(teacher)
  })
}