package teacher

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func update(g fiber.Router) {
  g.Put("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var draftMark bong.DraftMark
    json.Unmarshal(c.Body(), &draftMark)

    draftMark.Update()
  
    return c.JSON(draftMark)
  })

  g.Patch("/points/increase", authMiddleware, func (c *fiber.Ctx) error {
    points, err := bong.IncreasePoints(
      bson.M {
        "subjectID": c.Query("subjectID"),
        "studentID": c.Query("studentID"),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(points)
  })

  g.Patch("/points/decrease", authMiddleware, func (c *fiber.Ctx) error {
    points, err := bong.DecreasePoints(
      bson.M {
        "subjectID": c.Query("subjectID"),
        "studentID": c.Query("studentID"),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(points)
  })

  
}