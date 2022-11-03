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

    fmt.Println(grip.CheckTeacher("0723010405"))
    return c.JSON(teacher)
  })

  signup.Post("/basic", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    if (grip.CheckTeacher(body["phone"])) {
      return utils.MessageError(c, "există deja un cont cu numărul acesta")
    } else {
      teacher := grip.Teacher {
        LastName: body["lastName"],
        FirstName: body["firstName"],
      }
      teacher.Put()
      token := grip.GenTeacherToken(teacher)

      code := utils.GenCode()
      utils.SendSMS("+4" + body["phone"], code)

      hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
      grip.Set("code:" + body["phone"], string(hashedCode))

      if err != nil {
        return utils.MessageError(c, "problemă internă :(")
      }

      return c.JSON(map[string]interface{} {
        "teacher": teacher,
        "token": token,
      })
    }
  })

  signup.Post("/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := grip.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    if compareErr == nil {
      key := fmt.Sprintf("%v", c.Locals("key"))
      err := grip.UpdateTeacher(key, 
        base.Updates {
          "phone": body["phone"],
        },
      )
      if err != nil {
        return utils.MessageError(c, "problemă internă :(")
      }
    
      var teacher grip.Teacher
      utils.GetLocals(c, "teacher", &teacher)
      teacher.Phone = body["phone"]

      token := grip.GenTeacherToken(teacher)

      grip.Del("code:" + body["phone"])

      return c.JSON(map[string]interface{} {
        "teacher": teacher,
        "token": token,
      })
    } else {
      return utils.MessageError(c, "codul introdus este greșit")
    }
  })

  signup.Post("/password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    key := fmt.Sprintf("%v", c.Locals("key"))
    grip.UpdateTeacher(key, 
      base.Updates {
        "password": string(hashedPassword),
      },
    )

    teacher, _ := grip.GetTeacher(
      base.Query {
        {"key": c.Locals("key")},
      },
    )

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
    })
  })

  signup.Post("/passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(body["passcode"]), 10)
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    key := fmt.Sprintf("%v", c.Locals("key"))
    grip.UpdateTeacher(key, 
      base.Updates {
        "passcode": string(hashedPasscode),
      },
    )

    teacher, _ := grip.GetTeacher(
      base.Query {
        {"key": c.Locals("key")},
      },
    )

    return c.JSON(map[string]interface{} {
      "teacher": teacher,
    })
  })
}