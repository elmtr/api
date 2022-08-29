package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Encoding array
var Encoding string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
var CodeEncoding string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenID() string {
  var ID string
  for i := 0; i < 6; i++ {
    ID += string(Encoding[rand.Intn(64)])
  }

  return ID
}

func GenCode() string {
  var code string
  for i := 0; i < 4; i++ {
    code += string(CodeEncoding[rand.Intn(36)])
  }
  return code
}

func GetLocals(c *fiber.Ctx, name string, result interface{}) {
  json.Unmarshal([]byte(fmt.Sprintf("%v", c.Locals(name))), &result)
} 

func SetLocals(c *fiber.Ctx, name string,  data interface{}) {
	bytes, _ := json.Marshal(data)
	json := string(bytes)
	c.Locals(name, json)
}

func Error(c *fiber.Ctx, err error) error  {
  return c.Status(500).SendString(fmt.Sprintf("%v", err))
} 

func MessageError(c *fiber.Ctx, message string) error {
  return c.Status(401).JSON(bson.M{
    "message": message,
  })
}

