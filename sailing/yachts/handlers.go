package yachts

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reader func(interface{}, ...interface{}) *gorm.DB

type Writer func(interface{}) *gorm.DB

func Route(r *fiber.App, s *gorm.DB) {
	r.Get("/yachts", ListYachts(s.Find))
}

type Yacht struct {
	ID          uint
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Price       float32    `json:"price"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Updated     *time.Time `json:"updated"`
}

// @Summary List
// @Schemes
// @Description List yachts
// @Tags yachts
// @Accept application/json
// @Success 200 {object} []Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /yachts [get]
func ListYachts(find Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var yachts []Yacht

		if result := find(&yachts); result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(yachts)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}
