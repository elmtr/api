package teacher

import (
	"api/grip"
	"api/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func set(g fiber.Router) {
  g.Post("/marks", authMiddleware, func (c *fiber.Ctx) error {
    var mark grip.Mark
    json.Unmarshal(c.Body(), &mark)

    mark.Put()

    return c.JSON(mark)
  })

  g.Post("/truancies", authMiddleware, func (c *fiber.Ctx) error {
    var truancy grip.Truancy
    json.Unmarshal(c.Body(), &truancy)

    truancy.Put()

    return c.JSON(truancy)
  })

  g.Post("/draftmarks", authMiddleware, func (c *fiber.Ctx) error {
    var draftMark grip.DraftMark
    json.Unmarshal(c.Body(), &draftMark)

    draftMark.Put()

    return c.JSON(draftMark)
  })

  g.Post("/draftmarks/definitivate", authMiddleware, func (c *fiber.Ctx) error {
    key := fmt.Sprintf("%v", c.Query("key"))
    mark, err := grip.DefinitivateDraftMark(key)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(mark)
  })
}