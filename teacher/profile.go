package teacher

import (
	"api/grip"
	"api/utils"
	"encoding/json"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func profile(g fiber.Router) {
  profile := g.Group("/profile")

  profile.Get("/", authMiddleware, func (c *fiber.Ctx) error {
    var teacher grip.Teacher
    utils.GetLocals(c, "teacher", &teacher)

    return c.JSON(teacher)
  })

  profile.Post("/change-password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    teacher, err := grip.GetTeacher(
      base.Query {
        {"id": c.Locals("id")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPassword, err := bcrypt.GenerateFromPassword([]byte(body["newPassword"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = grip.UpdateTeacher(
      teacher.Key,
      base.Updates {
        "password": string(newPassword),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    teacher.Password = string(newPassword)

    return c.JSON(teacher)
  })

  profile.Post("/change-passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    teacher, err := grip.GetTeacher(
      base.Query {
        {"key": c.Locals("key")},
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(teacher.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPasscode, err := bcrypt.GenerateFromPassword([]byte(body["newPasscode"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = grip.UpdateTeacher(
      teacher.Key,
      base.Updates {
        "passcode": string(newPasscode),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    teacher.Passcode = string(newPasscode)

    return c.JSON(teacher)
  })
}