package main

import (
	"log"

	"api/admin"
	"api/bong"
	"api/parent"
	"api/student"
	"api/teacher"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var MongoURI string = "mongodb+srv://stevensun:stevensun@dev.ne1g1.mongodb.net/dev?retryWrites=true&w=majority"
var RedisOptions *redis.Options = &redis.Options {
  Addr: "127.0.0.1:6379", 
  Password: "",
  DB: 0,
}

func main() {
  app := fiber.New()

  bong.InitDB(MongoURI)
  bong.InitCache(RedisOptions)

  app.Get("/test", func (c *fiber.Ctx) error {
    return c.SendString("Hello there")
  })

  teacher.Routes(app)
  student.Routes(app)
  parent.Routes(app)
  admin.Routes(app)

  log.Fatal(app.Listen(":4200"))
}