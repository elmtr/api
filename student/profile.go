package student

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
    var student bong.Student
    utils.GetLocals(c, "student", &student)

    return c.JSON(student)
  })

  profile.Post("/change-password", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    student, err := bong.GetStudent(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPassword, err := bcrypt.GenerateFromPassword([]byte(body["newPassword"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = bong.UpdateStudent(
      student.ID,
      bson.M{
        "password": string(newPassword),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    student.Password = string(newPassword)

    return c.JSON(student)
  })

  profile.Post("/change-passcode", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    student, err := bong.GetStudent(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }
  
    compareErr := bcrypt.CompareHashAndPassword([]byte(student.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }
    
    newPasscode, err := bcrypt.GenerateFromPassword([]byte(body["newPasscode"]), 10)
    if err != nil {
      return utils.Error(c, err)
    }

    err = bong.UpdateStudent(
      student.ID,
      bson.M{
        "passcode": string(newPasscode),
      },
    )
    if err != nil {
      return utils.Error(c, err)
    }
    student.Passcode = string(newPasscode)

    return c.JSON(student)
  })
}