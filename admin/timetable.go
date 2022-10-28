package admin

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
    var admin grip.Admin
    utils.GetLocals(c, "admin", &admin)

    school, err := grip.GetSchool(admin.SchoolKey)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(school)
  })

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    periods, err := grip.GetPeriods(
      base.Query {
        {"subject.grade.key": c.Query("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })

  tt.Post("/", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    day, _ := strconv.Atoi(body["day"])
    interval, _ := strconv.Atoi(body["interval"])

    var subject grip.Subject
    var err error

    if (body["name"] == "Consiliere") {
      grade, err := grip.GetGrade(
        base.Query {
          {"key": body["gradeKey"]},
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
      subject = grip.Subject {
        Key: "0",
        Name: "Consiliere",
        Grade: grade,
      }
    } else {
      subject, err = grip.GetSubject(
        base.Query {
          {
            "name": body["name"],
            "grade.key": body["gradeKey"],
          },
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
    }

    period := grip.Period {
      Day: day,
      Interval: interval,
      Room: body["room"],
      Subject: subject,
    }
    period.Put()

    return c.JSON(period)
  })

  tt.Delete("/", authMiddleware, func (c *fiber.Ctx) error {
    key := c.Query("key")

    err := grip.DeletePeriod(key)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON("ok")
  })

  tt.Patch("/", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)
    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

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
      return utils.MessageError(c, "Nu s-a putut gasi materia introdusa")
    }

    err = grip.UpdatePeriod(c.Query("key"), subject, body["room"])
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(subject)
  })
}