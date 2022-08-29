package parent

import (
	"api/bong"
	"api/utils"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
  var token string

  authHeader := c.Get("Authorization")

  if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {
    token = strings.Fields(string(authHeader))[1]

    // we're just parsing the token
    claims := &bong.ParentClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return bong.JWTKey, nil
    })
    fmt.Println(claims)

    if err != nil {
      return utils.Error(c, err)
    }

    if !tkn.Valid {
      return utils.MessageError(c, "token not valid")
    }

    c.Locals("id", claims.ID)
    parent := bong.Parent {
      ID: claims.ID,
      FirstName: claims.FirstName,
      LastName: claims.LastName,
      Phone: claims.Phone,
      Students: claims.Students,
    }
    utils.SetLocals(c, "parent", parent)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}

func Routes(app *fiber.App) {
  g := app.Group("/v1/parent")

  signup(g)
  login(g)
}