package teacher

import (
	"api/grip"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func update(g fiber.Router) {
  g.Put("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var draftMark grip.DraftMark
    json.Unmarshal(c.Body(), &draftMark)

    draftMark.Update()
  
    return c.JSON(draftMark)
  })

  g.Patch("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    key := c.Query("key")

    err := grip.MotivateTruancy(key)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON("ok")
  })

  g.Patch("/points/increase", authMiddleware, func (c *fiber.Ctx) error {
    err := grip.IncreasePoints(c.Query("key"))
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON("ok")
  })

  g.Patch("/points/decrease", authMiddleware, func (c *fiber.Ctx) error {
    err := grip.DecreasePoints(c.Query("key"))
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON("ok")
  })  
}