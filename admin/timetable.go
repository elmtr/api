package admin

import (
	"api/bong"
	"api/utils"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func timetable(g fiber.Router) {
  tt := g.Group("/timetable")

  tt.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    periods, err := bong.GetPeriods(
      bson.M{
        "subject.grade.id": c.Query("id"),
      },
      bong.PeriodSort,
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(periods)
  })

  tt.Post("/", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])
    day, _ := strconv.Atoi(body["day"])
    interval, _ := strconv.Atoi(body["interval"])

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

    period := bong.Period {
      Day: day,
      Interval: interval,
      Room: body["room"],
      Subject: subject,
    }
    period.Insert()

    return c.JSON(period)
  })

  tt.Delete("/", authMiddleware, func (c *fiber.Ctx) error {
    id := c.Query("id")

    err := bong.DeletePeriod(
      bson.M{
        "id": id,
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON("ok")
  })

  tt.Patch("/", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)
    gradeNumber, _ := strconv.Atoi(body["gradeNumber"])

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

    period, err := bong.UpdatePeriod(c.Query("id"), subject, body["room"])
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(period)
  })
}