package teacher

import (
	"encoding/json"
	"fmt"

	"api/grip"
	"api/utils"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func signup(g fiber.Router) {
  signup := g.Group("/signup")

  signup.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)
    return c.JSON(teacher)
  })

  signup.Post("/basic", func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    json.Unmarshal(c.Body(), &teacher)

    teacher.Put()

    token := grip.GenTeacherToken(teacher)

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
      "token": token,
    })
  })

  signup.Post("/phone", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

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
    err := grip.UpdateTeacher(key, 
      base.Updates {
        "phone": body["phone"],
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
  
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)
    teacher.Phone = body["phone"]

    token := grip.GenTeacherToken(teacher)

    if compareErr == nil {
      grip.Del("code:" + body["phone"])
      return c.JSON(map[string]interface{} {
        "teacher": teacher,
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

    key := fmt.Sprintf("%V", c.Locals("key"))
    grip.UpdateTeacher(key, 
      base.Updates {
        "password": string(hashedPassword),
      },
    )

    if err != nil {
      return utils.Error(c, err)
    }

    teacher, _ := grip.GetTeacher(
      base.Query {
        {"id": c.Locals("id")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
    })
  })

  signup.Post("/passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(body["passcode"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    key := fmt.Sprintf("%v", c.Locals("key"))
    grip.UpdateTeacher(key, 
      base.Updates {
        "passcode": string(hashedPasscode),
      },
    )

    if err != nil {
      return utils.Error(c, err)
    }

    teacher, _ := grip.GetTeacher(
      base.Query {
        {"id": c.Locals("id")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
    })
  })
}