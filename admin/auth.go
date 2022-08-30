package admin

import (
	"api/bong"
	"api/utils"
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
    claims := &bong.AdminClaims{}
    tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface {}, error) {
      return bong.JWTKey, nil
    })

    if err != nil {
      return utils.Error(c, err)
    }

    if !tkn.Valid {
      return utils.MessageError(c, "token not valid")
    }

    c.Locals("id", claims.ID)
    admin := bong.Admin {
      ID: claims.ID,
      Email: claims.Email,
    }
    utils.SetLocals(c, "admin", admin)
  }

  if (token == "") {
    return utils.MessageError(c, "no token")
  }

  return c.Next()
}