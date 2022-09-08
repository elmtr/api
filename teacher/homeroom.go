package teacher

import (
	"api/bong"
	"api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func homeroom(g fiber.Router) {
  hr := g.Group("/homeroom")

  hr.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var teacher bong.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    return c.JSON(teacher)
  })

  hr.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher bong.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    students, err := bong.GetStudents(
      bson.M{
        "grade.id": teacher.Homeroom.ID,
      },
    )
    if err != nil {
      utils.Error(c, err)
    }

    return c.JSON(students)
  })

  hr.Get("/truancies", func (c *fiber.Ctx) error {
    id := c.Query("id")
    
    truancies, err := bong.GetTruancies(
      bson.M{
        "studentID": id,
      },
      bong.DateSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(truancies)
  })

  // TODO: also read the grade's timetable once i am done with figuring that out.
  // TODO: be able to close a year average mark: after the law gets clearer
}
