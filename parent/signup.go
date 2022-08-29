package parent

import (
	"api/bong"
	"api/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func signup(g fiber.Router) {
  signup := g.Group("/signup")

  signup.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var parent bong.Parent
    utils.GetLocals(c.Locals("parent"), &parent)
    return c.JSON(parent)
  })

  signup.Post("/basic", func (c *fiber.Ctx) error {
    var parent bong.Parent
    json.Unmarshal(c.Body(), &parent)

    parent.Insert()

    token, err := bong.GenParentToken(parent)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "parent": parent,
      "token": token,
    })
  })

  signup.Post("/phone", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    err := bong.UpdateParent(c.Locals("id"), bson.M{"phone": body["phone"]})

    if err != nil {
      return utils.Error(c, err)
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

  signup.Post("/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := bong.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    var parent bong.Parent
    utils.GetLocals(c.Locals("parent"), &parent)
    fmt.Println(parent)
    parent.Phone = body["phone"]

    token, err := bong.GenParentToken(parent)
    if err != nil {
      return utils.Error(c, err)
    }

    if compareErr == nil {
      bong.Del("code:" + body["phone"])
      return c.JSON(bson.M{
        "token": token,
      })
    } else {
      return utils.MessageError(c, "Codul introdus este greșit")
    }
  })

  signup.Post("/password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    bong.UpdateParent(c.Locals("id"), bson.M{"password": string(hashedPassword)})

    if err != nil {
      return utils.Error(c, err)
    }

    parent, _ := bong.GetParent(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "parent": parent,
    })
  })

  signup.Post("/passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(body["passcode"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    bong.UpdateParent(c.Locals("id"), bson.M{"passcode": string(hashedPasscode)})

    if err != nil {
      return utils.Error(c, err)
    }

    parent, _ := bong.GetParent(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "parent": parent,
    })
  })
}