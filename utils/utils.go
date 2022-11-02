package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Encoding array
var Encoding string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
var CodeEncoding string = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenID() string {
  var ID string
  for i := 0; i < 6; i++ {
    ID += string(Encoding[rand.Intn(64)])
  }

  return ID
}

func GenCode() string {
  var code string
  for i := 0; i < 6; i++ {
    code += strconv.Itoa(rand.Intn(10));
  }
  fmt.Println(code)
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
  return c.Status(401).JSON(map[string]interface{} {
    "message": message,
  })
}

