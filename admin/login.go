package admin

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func login(g fiber.Router) {
  login := g.Group("/login")

  login.Post("/", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    admin, err := bong.GetAdmin(body["email"])
    if err != nil {
      return utils.Error(c, err)
    }

    compareErr := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "parola e gresita")
    }

    token, err := bong.GenAdminToken(admin)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "admin": admin,
      "token": token,
    })
  })
}