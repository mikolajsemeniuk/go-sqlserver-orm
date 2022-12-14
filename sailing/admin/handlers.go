package admin

import (
	"encoding/json"
	"server/sailing/security"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reader func(query interface{}, args ...interface{}) (tx *gorm.DB)

type Writer func(value interface{}) (tx *gorm.DB)

type Remover func(value interface{}, conds ...interface{}) (tx *gorm.DB)

func Route(r *fiber.App, s *gorm.DB) {
	r.Get("/admin/list", security.Authorize("admin"), ListYachts(s.Find))
	r.Post("/admin/create", security.Authorize("admin"), CreateYacht(s.Create))
	r.Delete("/admin/remove/:id", security.Authorize("admin"), RemoveYacht(s.Delete))
}

// @Summary Create
// @Schemes
// @Description Create new yacht
// @Tags admin
// @Accept application/json
// @Success 200 {object} []Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/list [get]
func ListYachts(r Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var yachts []Yacht

		if result := r(&yachts); result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(yachts)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}

// @Summary Create
// @Schemes
// @Description Create new yacht
// @Tags admin
// @Accept application/json
// @Param payload body CreateYachtRequest true "body"
// @Success 200 {object} Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/create [post]
func CreateYacht(w Writer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request CreateYachtRequest

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		yacht := Yacht{
			Name:        request.Name,
			Type:        request.Type,
			Price:       request.Price,
			Image:       request.Image,
			Description: request.Description,
			Created:     time.Now(),
		}
		if result := w(&yacht); result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(yacht)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}

// @Summary Remove
// @Schemes
// @Description Remove yacht
// @Tags admin
// @Accept application/json
// @Param id path string true "id"
// @Success 200
// @Success 400
// @Success 404
// @Failure 500
// @Failure 503
// @Router /admin/remove/{id} [delete]
func RemoveYacht(r Remover) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		yacht := Yacht{
			ID: uint(id),
		}

		result := r(&yacht)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		if result.RowsAffected == 0 {
			return fiber.NewError(fiber.StatusNotFound, "not found")
		}

		response, err := json.Marshal(yacht)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}
