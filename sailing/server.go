package sailing

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret = "some secret"
var cookie = "cookie-name"

type Reader func(query interface{}, args ...interface{}) (tx *gorm.DB)

type Writer func(value interface{}) (tx *gorm.DB)

// SignIn
//
// SignUp is used for:
//   - validating data
//   - checking if user with given email exists and if yes return error
//   - generating password hash
//   - write user to database
//   - creating token
func SignUp(r Reader, w Writer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request SignUpRequest

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		result := r("email = ?", request.Email).Find(&Account{})
		if result.Error != nil {
			return fiber.NewError(fiber.StatusBadRequest, result.Error.Error())
		}

		if result.RowsAffected != 0 {
			return fiber.NewError(fiber.StatusBadRequest, "user with this email already exists")
		}

		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
		if err != nil {
			return err
		}

		entity := Account{
			Name:     request.Name,
			Email:    request.Email,
			Password: password,
			Created:  time.Now(),
		}
		if w(&entity).Error != nil {
			return fiber.NewError(fiber.StatusBadRequest, result.Error.Error())
		}

		claims := jwt.MapClaims{"Issuer": entity.ID, "ExpiresAt": time.Now().Add(time.Hour).Unix()}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		value, err := token.SignedString([]byte(secret))
		if err != nil {
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:     "cookie-name",
			Value:    value,
			Secure:   true,
			HTTPOnly: false,
			SameSite: "none",
			Expires:  time.Now().Add(time.Duration(1) * time.Hour),
		})

		return c.SendString("SignUp")
	}
}

func SignIn(r Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request SignInRequest

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var account Account
		result := r("email = ?", request.Email).Find(&account)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusBadRequest, result.Error.Error())
		}

		if result.RowsAffected != 1 {
			return fiber.NewError(fiber.StatusBadRequest, "no user with this email address")
		}

		if err := bcrypt.CompareHashAndPassword(account.Password, []byte(request.Password)); err != nil {
			return fiber.NewError(http.StatusBadRequest, "incorrect password")
		}

		claims := jwt.MapClaims{"Issuer": account.ID, "ExpiresAt": time.Now().Add(time.Hour).Unix()}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		value, err := token.SignedString([]byte(secret))
		if err != nil {
			return err
		}

		c.Cookie(&fiber.Cookie{
			Name:     cookie,
			Value:    value,
			Secure:   true,
			HTTPOnly: false,
			SameSite: "none",
			Expires:  time.Now().Add(time.Duration(1) * time.Hour),
		})

		return c.SendString("SignIn")
	}
}

// @Summary Logout
// @Schemes
// @Description Logout existing user
// @Tags account
// @Accept application/json
// @Success 200 {object} string
// @Router /account/logout [get]
func SignOut() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     cookie,
			Value:    "",
			Secure:   true,
			HTTPOnly: false,
			SameSite: "none",
			Expires:  time.Now().Add(-time.Second),
		})
		return c.SendString("user logout")
	}
}

// @Summary Authorize
// @Schemes
// @Description Authorize existing user
// @Tags account
// @Accept application/json
// @Success 200 {object} string
// @Success 401 {object} string
// @Router /account/authorize [get]
func Authorize() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies(cookie)
		parser := func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil }

		token, err := jwt.Parse(cookie, parser)
		if err != nil {
			return fiber.NewError(http.StatusUnauthorized, "unauthorized")
		}

		payload := token.Claims.(jwt.MapClaims)

		id := payload["Issuer"].(float64)

		return c.SendString(fmt.Sprintf("user with id: %d", int(id)))
	}
}
