package teacher

import (
	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    subjectsKeyList := []string {}
    for _, subject := range teacher.Subjects {
      subjectsKeyList = append(subjectsKeyList, subject.Key)
    }

    // TODO: do this operation for each and every subject 
    // in the subjectsKeyList: deta doesn't support anything similar to $in
    periods, err := grip.GetPeriods(
      base.Query {
        {"subject.key?contains": subjectsKeyList},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })
}