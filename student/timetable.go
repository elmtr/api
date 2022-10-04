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

  tt.Get("/school", authMiddleware, func (c *fiber.Ctx) error {
    var student grip.Student
    utils.GetLocals(c, "student", &student)

    school, err := grip.GetSchool(student.Grade.SchoolKey)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(school)
  })

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

    var student grip.Student
    utils.GetLocals(c, "student", &student)

    interval, _ := strconv.Atoi(c.Query("interval"))
    day, _ := strconv.Atoi(c.Query("day"))

    currentPeriod, err := grip.GetPeriod(
      base.Query {
        {
          "interval": interval,
          "day": day,
          "subject.grade.key": student.Grade.Key,
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
          "subject.grade.key": student.Grade.Key,
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