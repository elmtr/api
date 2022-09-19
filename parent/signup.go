package parent

import (
	"api/grip"
	"api/utils"
	"encoding/json"
	"fmt"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func signup(g fiber.Router) {
  signup := g.Group("/signup")

  signup.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var parent grip.Parent
    utils.GetLocals(c, "parent", &parent)
    return c.JSON(parent)
  })

  signup.Post("/basic", func (c *fiber.Ctx) error {
    var parent grip.Parent
    json.Unmarshal(c.Body(), &parent)

    parent.Put()

    token := grip.GenParentToken(parent)

    return c.JSON(map[string]interface{} {
      "parent": parent,
      "token": token,
    })
  })

  signup.Post("/phone", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    key := fmt.Sprintf("%v", c.Locals("id"))
    err := grip.UpdateParent(key, 
      base.Updates {
        "phone": body["phone"],
      },
    )

    if err != nil {
      return utils.Error(c, err)
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

  signup.Post("/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := grip.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    key := fmt.Sprintf("%v", c.Locals("key"))
    err := grip.UpdateParent(key, 
      base.Updates {
        "phone": body["phone"],
      },
    ) 

    if err != nil {
      return utils.Error(c, err)
    }

    var parent grip.Parent
    utils.GetLocals(c, "parent", &parent)
    parent.Phone = body["phone"]

    token := grip.GenParentToken(parent)

    if compareErr == nil {
      grip.Del("code:" + body["phone"])
      return c.JSON(map[string]interface{} {
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

    key := fmt.Sprintf("%v", c.Locals("id"))
    grip.UpdateParent(key, base.Updates {
        "password": string(hashedPassword),
      },
    )

    if err != nil {
      return utils.Error(c, err)
    }

    parent, _ := grip.GetParent(
      base.Query {
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(map[string]interface{} {
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

    key := fmt.Sprintf("%v", c.Locals("id"))
    grip.UpdateParent(key, base.Updates {
        "passcode": string(hashedPasscode),
      },
    )

    if err != nil {
      return utils.Error(c, err)
    }

    parent, _ := grip.GetParent(
      base.Query {
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(map[string]interface{} {
      "parent": parent,
    })
  })
}