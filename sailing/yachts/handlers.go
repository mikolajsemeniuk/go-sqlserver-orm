package yachts

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reader func(interface{}, ...interface{}) *gorm.DB

type Writer func(interface{}) *gorm.DB

func Route(r *fiber.App, s *gorm.DB) {
	r.Get("/yachts", ListYachts(s.Find))
	r.Get("/yachts/:id", FindYacht(s.Find))
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

// @Summary Find
// @Schemes
// @Description Find yacht
// @Tags yachts
// @Accept application/json
// @Param id path string true "id"
// @Success 200 {object} Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /yachts/{id} [get]
func FindYacht(first Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		yacht := Yacht{
			ID: uint(id),
		}

		result := first(&yacht, id)
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return fiber.NewError(fiber.StatusNotFound, `{ "message": "not found" }`)
		case nil:
		default:
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(yacht)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}
