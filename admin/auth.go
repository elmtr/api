package admin

import (
	"api/grip"
	"api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    admin, err := grip.ParseAdminToken(token)
    if err != nil {
      return utils.Error(c, err)
    }

    c.Locals("key", admin.Key)
    utils.SetLocals(c, "admin", admin)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}