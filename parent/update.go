package parent

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func update(g fiber.Router) {
  g.Post("/students", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    _, err := bong.GetStudent(
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

  g.Post("/students/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := bong.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    if compareErr == nil {
      // bong.Del("code:" + body["phone"])

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
      parentStudent := bong.ParentStudent {
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