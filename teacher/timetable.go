package teacher

import (
	"api/bong"
	"api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher bong.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    subjectsIDList := []string {}
    for _, subject := range teacher.Subjects {
      subjectsIDList = append(subjectsIDList, subject.ID)
    }

    periods, err := bong.GetPeriods(
      bson.M{
        "subject.id": bson.M{
          "$in": subjectsIDList,
        },
      },
      bong.PeriodSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })
}