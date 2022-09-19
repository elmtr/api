package teacher

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

    teacher, err := grip.ParseTeacherToken(token)

    if err != nil {
      return utils.Error(c, err)
    }
    
    c.Locals("key", teacher.Key)
    utils.SetLocals(c, "teacher", teacher)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}