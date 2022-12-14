package sailing

import (
	"encoding/json"
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

// @Summary Register
// @Schemes
// @Description Register new account
// @Tags account
// @Accept application/json
// @Param payload body RegisterRequest true "body"
// @Success 200 {object} Account
// @Success 400
// @Failure 500
// @Failure 503
// @Router /account/register [post]
func Register(r Reader, w Writer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request RegisterRequest

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		result := r("email = ?", request.Email).Find(&Account{})
		if result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		if result.RowsAffected != 0 {
			return fiber.NewError(fiber.StatusBadRequest, "user with this email already exists")
		}

		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		account := Account{
			Name:     request.Name,
			Email:    request.Email,
			Password: password,
			Role:     "user",
			Created:  time.Now(),
		}

		if w(&account).Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		claims := jwt.MapClaims{"Issuer": account.ID, "Role": account.Role, "ExpiresAt": time.Now().Add(time.Hour).Unix()}
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

		response, _ := json.Marshal(account)
		return c.Send(response)
	}
}

// @Summary Login
// @Schemes
// @Description Login existing user
// @Tags account
// @Accept application/json
// @Param payload body LoginRequest true "body"
// @Success 200 {object} Account
// @Success 400
// @Failure 500
// @Failure 503
// @Router /account/login [post]
func Login(r Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request LoginRequest

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

		claims := jwt.MapClaims{"Issuer": account.ID, "Role": account.Role, "ExpiresAt": time.Now().Add(time.Hour).Unix()}
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

		response, _ := json.Marshal(account)
		return c.Send(response)
	}
}

// @Summary Logout
// @Schemes
// @Description Logout existing user
// @Tags account
// @Accept application/json
// @Success 200 {object} string
// @Router /account/logout [get]
func Logout() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     cookie,
			Value:    "",
			Secure:   true,
			HTTPOnly: false,
			SameSite: "none",
			Expires:  time.Now().Add(-time.Second),
		})
		return c.Send([]byte(`{ }`))
	}
}
