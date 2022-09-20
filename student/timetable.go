package student

import (
	"api/grip"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    periods, err := grip.GetPeriods(
      base.Query {
        {"subject.grade.key": student.Grade.Key},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })

  tt.Get("/widget", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    interval, _ := strconv.Atoi(c.Query("interval"))
    day, _ := strconv.Atoi(c.Query("day"))
    gradeKey := c.Query("gradeKey")

    currentPeriod, err := grip.GetPeriod(
      base.Query {
        {
          "interval": interval,
          "day": day,
          "subject.grade.key": gradeKey,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    nextPeriod, err := grip.GetPeriod(
      base.Query {
        {
          "interval": interval + 1,
          "day": day,
          "subject.grade.key": gradeKey,
        },
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    return c.JSON(map[string]interface{} {
      "currentPeriod": currentPeriod,
      "nextPeriod": nextPeriod,
    })
  })
}