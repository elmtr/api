package parent

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

  login.Get("/test", func (c *fiber.Ctx) error {
    return c.SendString("hello")
  })

  login.Post("", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := bong.GetParent(bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdusă nu este validă")
    }
    
    code := utils.GenCode()
    utils.SendSMS("+4" + body["phone"], code)

    hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
    bong.Set("code:" + body["phone"], string(hashedCode))

    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(body["phone"])
  })

  login.Post("/verify-code", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := bong.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))
  
    parent, err := bong.GetParent(bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }

    if compareErr == nil {
      bong.Del("code:" + body["phone"])
      return c.JSON(bson.M{
        "parent": parent,
      })
    } else {
      return utils.MessageError(c, "Codul introdus este greșit")
    }
  })

  login.Post("/update", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := bong.GetParent(bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }

    token, err := bong.GenParentToken(parent)
    if err != nil {
      return utils.Error(c, err)
    }

    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }

    return c.JSON(bson.M{
      "parent": parent,
      "token": token,
    })
  })
}