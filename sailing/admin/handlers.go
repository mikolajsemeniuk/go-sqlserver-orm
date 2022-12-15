package admin

import (
	"encoding/json"
	"server/sailing/security"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reader func(interface{}, ...interface{}) *gorm.DB

type Writer func(interface{}) *gorm.DB

func Route(r *fiber.App, s *gorm.DB) {
	r.Get("/admin/list", security.Authorize("admin"), ListYachts(s.Find))
	r.Get("/admin/yacht/:id", security.Authorize("admin"), FindYacht(s.First))
	r.Post("/admin/create", security.Authorize("admin"), CreateYacht(s.Create))
	r.Patch("/admin/yacht/update", security.Authorize("admin"), UpdateYacht(s.First, s.Save))
	r.Delete("/admin/remove/:id", security.Authorize("admin"), RemoveYacht(s.Delete))
}

// @Summary Create
// @Schemes
// @Description List yachts
// @Tags admin
// @Accept application/json
// @Success 200 {object} []Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/list [get]
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
// @Tags admin
// @Accept application/json
// @Param id path string true "id"
// @Success 200 {object} Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/yacht/{id} [get]
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

// @Summary Create
// @Schemes
// @Description Create yacht
// @Tags admin
// @Accept application/json
// @Param payload body CreateYachtRequest true "body"
// @Success 200 {object} Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/create [post]
func CreateYacht(create Writer) func(c *fiber.Ctx) error {
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

		if result := create(&yacht); result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(yacht)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}

// @Summary Update
// @Schemes
// @Description Update yacht
// @Tags admin
// @Accept application/json
// @Param payload body UpdateYachtRequest true "body"
// @Success 200 {object} Yacht
// @Success 400
// @Failure 500
// @Failure 503
// @Router /admin/yacht/update [patch]
func UpdateYacht(first Reader, update Writer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var request UpdateYachtRequest

		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var yacht Yacht

		result := first(&yacht, request.ID)
		switch result.Error {
		case gorm.ErrRecordNotFound:
			return fiber.NewError(fiber.StatusNotFound, `{ "message": "not found" }`)
		case nil:
		default:
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		updated := time.Now()
		yacht.Name = request.Name
		yacht.Type = request.Type
		yacht.Price = request.Price
		yacht.Image = request.Image
		yacht.Description = request.Description
		yacht.Updated = &updated

		if err := update(&yacht).Error; err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
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
func RemoveYacht(remove Reader) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		yacht := Yacht{
			ID: uint(id),
		}

		result := remove(&yacht)
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
