package parent

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func profile(g fiber.Router) {
  profile := g.Group("/profile")

  profile.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var parent bong.Parent
    utils.GetLocals(c, "parent", &parent)

    return c.JSON(parent)
  })

  profile.Post("/change-password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    parent, err := bong.GetParent(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(parent.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPassword, err := bcrypt.GenerateFromPassword([]byte(body["newPassword"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = bong.UpdateParent(
      parent.ID,
      bson.M{
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

    parent, err := bong.GetParent(bson.M{"id": c.Locals("id")})
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

    err = bong.UpdateParent(
      parent.ID,
      bson.M{
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

    _, err := bong.GetParent(
      bson.M{
        "phone": body["phone"],
      },
    )
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

  profile.Post("/students/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := bong.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    if compareErr == nil {
      bong.Del("code:" + body["phone"])

      var parent bong.Parent
      utils.GetLocals(c, "parent", &parent)

      student, err := bong.GetStudent(
        bson.M{
          "phone": body["phone"],
        },
      )
      if err != nil {
        return utils.Error(c, err)
      }
      parentStudent := bong.ParentStudent{
        ID: student.ID,
        FirstName: student.FirstName,
        LastName: student.LastName,
      }

      parent.Students, err = bong.AddParentStudent(parent.ID, parent.Students, parentStudent)
      if err != nil {
        return utils.Error(c, err)
      }
      
      token, err := bong.GenParentToken(parent)
      if err != nil {
        return utils.Error(c, err)
      }

      return c.JSON(bson.M{ 
          "token": token,
          "parent": parent,
      })
    } else {
      return utils.MessageError(c, "Codul introdus este greșit")
    }
  })
}