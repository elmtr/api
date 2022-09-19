package student

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

    student, err := grip.ParseStudentToken(token)
    if err != nil {
      return utils.Error(c, err)
    }

    c.Locals("key", student.Key)
    utils.SetLocals(c, "student", student)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}