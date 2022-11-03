package teacher

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

  login.Get("/test", func (c *fiber.Ctx) error {
    return c.SendString("hello")
  })

  login.Post("", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    teacher, err := grip.GetTeacher(
      base.Query {
        {"phone": body["phone"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "parola introdusă nu este corectă")
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
  
    teacher, err := grip.GetTeacher(
      base.Query{
        {"phone": body["phone"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    if compareErr == nil {
      grip.Del("code:" + body["phone"])
      return c.JSON(map[string]interface{} {
        "teacher": teacher,
      })
    } else {
      return utils.MessageError(c, "codul introdus este greșit")
    }
  })

  login.Post("/update", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    teacher, err := grip.GetTeacher(
      base.Query {
        {"phone": body["phone"]},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    token:= grip.GenTeacherToken(teacher)

    compareErr := bcrypt.CompareHashAndPassword([]byte(teacher.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "PIN-ul introdus nu este corect")
    }

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
      "token": token,
    })
  })
}