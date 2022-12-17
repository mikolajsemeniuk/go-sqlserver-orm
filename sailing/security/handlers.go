package security

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var secret = "some secret"
var cookie = "cookie-name"

const ID = "id"
const Role = "role"

func Authorize(role string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies(cookie)
		parser := func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil }

		token, err := jwt.Parse(cookie, parser)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, "only logged in users can perform this operation")
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(http.StatusUnauthorized, "invalid claims")
		}

		actual, ok := payload["Role"].(string)
		if !ok {
			return fiber.NewError(http.StatusUnauthorized, "missing role")
		}

		if role != actual {
			return fiber.NewError(http.StatusForbidden, "you cannot perform this operation")
		}

		c.Locals("id", payload["Issuer"])
		return c.Next()
	}
}
