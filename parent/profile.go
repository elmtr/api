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

func profile(g fiber.Router) {
  profile := g.Group("/profile")

  profile.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var parent grip.Parent
    utils.GetLocals(c, "parent", &parent)

    return c.JSON(parent)
  })

  profile.Post("/change-password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := grip.GetParent(
      base.Query {
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    fmt.Println(parent.Key)
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPassword, err := bcrypt.GenerateFromPassword([]byte(body["newPassword"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    key := fmt.Sprintf("%v", c.Locals("key"))
    err = grip.UpdateParent(
      key,
      base.Updates {
        "password": string(newPassword),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    parent.Password = string(newPassword)

    return c.JSON(parent)
  })

  profile.Post("/change-passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := grip.GetParent(
      base.Query {
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPasscode, err := bcrypt.GenerateFromPassword([]byte(body["newPasscode"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = grip.UpdateParent(
      parent.Key,
      base.Updates {
        "passcode": string(newPasscode),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    parent.Passcode = string(newPasscode)

    return c.JSON(parent)
  })

  profile.Post("/students", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    _, err := grip.GetParent(
      base.Query {
        {"phone": body["phone"]},
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

  profile.Post("/students/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := grip.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    if compareErr == nil {
      grip.Del("code:" + body["phone"])

      var parent grip.Parent
      utils.GetLocals(c, "parent", &parent)

      student, err := grip.GetStudent(
        base.Query {
          {"phone": body["phone"]},
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
      parentStudent := grip.ParentStudent{
        ID: student.Key,
        FirstName: student.FirstName,
        LastName: student.LastName,
      }

      parent.Students, err = grip.AddParentStudent(parent.Key, parent.Students, parentStudent)
      if err != nil {
        return utils.Error(c, err)
      }
      
      token := grip.GenParentToken(parent)

      return c.JSON(map[string]interface{} { 
          "token": token,
          "parent": parent,
      })
    } else {
      return utils.MessageError(c, "Codul introdus este greșit")
    }
  })
}