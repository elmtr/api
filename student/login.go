package student

import (
	"api/bong"
	"api/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

    student, err := bong.GetStudent(bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }

    compareErr := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body["password"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }

    return c.JSON(bson.M{
      "student": student,
    })
  })

  login.Post("/update", func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    student, err := bong.GetStudent(bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }

    token, err := bong.GenStudentToken(student)
    if err != nil {
      return utils.Error(c, err)
    }

    compareErr := bcrypt.CompareHashAndPassword([]byte(student.Passcode), []byte(body["passcode"]))
    if compareErr != nil {
      return utils.MessageError(c, "Parola introdus nu este valid")
    }

    return c.JSON(bson.M{
      "student": student,
      "token": token,
    })
  })
}