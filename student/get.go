package student

import (
	"api/bong"
	"api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func get(g fiber.Router) {
  g.Get("/marks", authMiddleware, func (c *fiber.Ctx) error {
    var student bong.Student
    utils.GetLocals(c, "student", &student)

    marks, err := bong.GetMarks(
      bson.M{
        "subjectID": c.Query("subjectID"),
        "studentID": student.ID,
      },
      bong.DateSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(marks)
  })

  g.Get("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    var student bong.Student
    utils.GetLocals(c, "student", &student)

    truancies, err := bong.GetTruancies(
      bson.M{
        "subjectID": c.Query("subjectID"),
        "studentID": student.ID,
      },
      bong.DateSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(truancies)
  })

  g.Get("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var student bong.Student
    utils.GetLocals(c, "student", &student)

    draftMarks, err := bong.GetDraftMarks(
      bson.M{
        "subjectID": c.Query("subjectID"),
        "studentID": student.ID,
      },
      bong.DateSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(draftMarks)
  })

  g.Get("/points", authMiddleware, func (c *fiber.Ctx) error {
    var student bong.Student
    utils.GetLocals(c, "student", &student)

    points, err := bong.GetPoints(
      bson.M{
        "subjectID": c.Query("subjectID"),
        "studentID": student.ID,
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(points)
  })
}