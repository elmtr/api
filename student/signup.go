package student

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
      var student grip.Student
      utils.GetLocals(c, "student", &student)
      return c.JSON(student)
  })
  
  signup.Post("/basic", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    if (grip.CheckStudent(body["phone"])) {
      return utils.MessageError(c, "există deja un cont cu numarul acesta")
    } else {
      student := grip.Student {
        LastName: body["lastName"],
        FirstName: body["firstName"],
      }
      student.Put()
      token := grip.GenStudentToken(student)

      code := utils.GenCode()
      utils.SendSMS("+4" + body["phone"], code)

      hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
      grip.Set("code:" + body["phone"], string(hashedCode))

      if err != nil {
        return utils.MessageError(c, "problemă internă :(")
      }

      return c.JSON(map[string]interface{} {
        "student": student,
        "token": token,
      })
    }
  })

  signup.Post("/phone", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    code := utils.GenCode()
    utils.SendSMS("+4" + body["phone"], code)

    hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
    grip.Set("code:" + body["phone"], string(hashedCode))

    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    return c.JSON(body["phone"])
  })

  signup.Post("/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := grip.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    key := fmt.Sprintf("%v", c.Locals("key"))
    err := grip.UpdateStudent(key, 
      base.Updates{
        "phone": body["phone"],
      },
    )
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }
  
    var student grip.Student
    utils.GetLocals(c, "student", &student)
    student.Phone = body["phone"]

    token := grip.GenStudentToken(student)

    if compareErr == nil {
      grip.Del("code:" + body["phone"])
      return c.JSON(map[string]interface{} {
        "token": token,
        "student": student,
      })
    } else {
      return utils.MessageError(c, "codul introdus este greșit")
    }
  })

  signup.Post("/grade", authMiddleware, func (c *fiber.Ctx) error {
    var grade grip.Grade
    json.Unmarshal(c.Body(), &grade)

    key := fmt.Sprintf("%v", c.Locals("key"))
    student, err := grip.StudentSetup(key, grade)
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    token := grip.GenStudentToken(student)

    return c.JSON(map[string]interface{} {
      "student": student,
      "token": token,
    })
  })

  signup.Post("/password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body["password"]), 10)
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    key := fmt.Sprintf("%v", c.Locals("key"))
    grip.UpdateStudent(key, 
      base.Updates {
        "password": string(hashedPassword),
      },
    )

    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    student, _ := grip.GetStudent(
      base.Query{
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.MessageError(c, "problemă internă :(")
    }

    return c.JSON(map[string]interface{} {
      "student": student,
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
    grip.UpdateStudent(key, 
      base.Updates {
        "passcode": string(hashedPasscode),
      },
    )

    student, _ := grip.GetStudent(
      base.Query {
        {"key": c.Locals("key")},
      },
    )

    return c.JSON(map[string]interface{} {
      "student": student,
    })
  })
}