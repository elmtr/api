package teacher

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func set(g fiber.Router) {
  g.Post("/marks", authMiddleware, func (c *fiber.Ctx) error {
    var mark bong.Mark
    json.Unmarshal(c.Body(), &mark)

    mark.Insert()

    return c.JSON(mark)
  })

  g.Post("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    var truancy bong.Truancy
    json.Unmarshal(c.Body(), &truancy)

    truancy.Insert()

    return c.JSON(truancy)
  })

  g.Post("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var draftMark bong.DraftMark
    json.Unmarshal(c.Body(), &draftMark)

    draftMark.Insert()

    return c.JSON(draftMark)
  })

  g.Post("/draftmarks/definitivate", authMiddleware, func (c *fiber.Ctx) error {
    mark, err := bong.DefinitivateDraftMark(
      bson.M{
        "id": c.Query("id"),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(mark)
  })
}