package admin

import (
	"api/grip"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func login(g fiber.Router) {
  login := g.Group("/login")

  login.Post("/", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    admin, err := grip.GetAdmin(body["email"])
    if err != nil {
      return utils.Error(c, err)
    }

    compareErr := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "parola e gresita")
    }

    token := grip.GenAdminToken(admin)

    return c.JSON(map[string]interface{} {
      "admin": admin,
      "token": token,
    })
  })

  login.Post("/start", func (c *fiber.Ctx) error {
    admin := grip.Admin {}
    admin.Put()
    return c.SendString("ok")
  })
}

