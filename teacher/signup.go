package teacher

import (
	"encoding/json"

	"api/bong"
	"api/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func signup(g fiber.Router) {
  signup := g.Group("/signup")

  signup.Get("/test", authMiddleware, func (c *fiber.Ctx) error {
    var teacher bong.Teacher
    utils.GetLocals(c, "teacher", &teacher)
    return c.JSON(teacher)
  })

  signup.Post("/basic", func (c *fiber.Ctx) error {
    var teacher bong.Teacher
    json.Unmarshal(c.Body(), &teacher)

    teacher.Insert()

    token, err := bong.GenTeacherToken(teacher)
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
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
    bong.Set("code:" + body["phone"], string(hashedCode))

    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(body["phone"])
  })

  signup.Post("/verify-code", authMiddleware, func (c *fiber.Ctx) error {
    var body map[string]string
    json.Unmarshal(c.Body(), &body)

    hashedCode, _ := bong.Get("code:" + body["phone"])

    compareErr := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(body["code"]))

    err := bong.UpdateTeacher(c.Locals("id"), bson.M{"phone": body["phone"]})
    if err != nil {
      return utils.Error(c, err)
    }
  
    var teacher bong.Teacher
    utils.GetLocals(c, "teacher", &teacher)
    teacher.Phone = body["phone"]

    token, err := bong.GenTeacherToken(teacher)
    if err != nil {
      return utils.Error(c, err)
    }

    if compareErr == nil {
      bong.Del("code:" + body["phone"])
      return c.JSON(bson.M{
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

    bong.UpdateTeacher(c.Locals("id"), bson.M{"password": string(hashedPassword)})

    if err != nil {
      return utils.Error(c, err)
    }

    teacher, _ := bong.GetTeacher(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
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

    bong.UpdateTeacher(c.Locals("id"), bson.M{"passcode": string(hashedPasscode)})

    if err != nil {
      return utils.Error(c, err)
    }

    teacher, _ := bong.GetTeacher(bson.M{"id": c.Locals("id")})
    if err != nil {
      return utils.Error(c, err)
    }

    return c.JSON(bson.M{
      "teacher": teacher,
    })
  })
}