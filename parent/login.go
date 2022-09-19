package parent

import (
	"api/grip"
	"api/utils"
	"encoding/json"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func login(g fiber.Router) {
  login := g.Group("/login")

  login.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var parent grip.Parent
    utils.GetLocals(c, "parent", &parent)

    return c.JSON(parent)
  })

  login.Post("", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := grip.GetParent(
      base.Query{
        {"phone": body["phone"]},
      },
    )
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
    grip.Set("code:" + body["phone"], string(hashedCode))

    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(body["phone"])
  })

  login.Post("/verify-code", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := grip.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))
  
    parent, err := grip.GetParent(
      base.Query {
        {"phone": body["phone"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    if compareErr == nil {
      grip.Del("code:" + body["phone"])
      return c.JSON(map[string]interface{} {
        "parent": parent,
      })
    } else {
      return utils.MessageError(c, "Codul introdus este greșit")
    }
  })

  login.Post("/update", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := grip.GetParent(
      base.Query {
        {"phone": body["phone"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    token := grip.GenParentToken(parent)

    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }

    return c.JSON(map[string]interface{} {
      "parent": parent,
      "token": token,
    })
  })
}