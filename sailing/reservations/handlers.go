package reservations

import (
	"encoding/json"
	"server/sailing/security"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Writer func(interface{}) *gorm.DB

type Reader func(interface{}, ...interface{}) *gorm.DB

func Route(r *fiber.App, s *gorm.DB) {
	r.Post("/reservations", security.Authorize("user"), Create(s, s.Create))
}

// @Summary Create
// @Schemes
// @Description Create reservation
// @Tags reservations
// @Accept application/json
// @Param payload body Request true "body"
// @Success 200 {object} Reservation
// @Success 400
// @Failure 500
// @Failure 503
// @Router /reservations [post]
func Create(s *gorm.DB, create Writer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		accountID, ok := c.Locals("id").(float64)
		if !ok {
			return fiber.NewError(fiber.StatusInternalServerError, `{ "message": "account id is missing or invalid"}`)
		}

		var request Request
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if request.From.After(request.To) {
			return fiber.NewError(fiber.StatusBadRequest, `{ "message": "from cannot be after to" }`)
		}

		var reservation Reservation
		query := `[from] between ? and ? or [to] between ? and ?`
		result := s.Table("reservations").Where(query, request.From, request.To, request.From, request.To).Find(&reservation)
		if result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		if result.RowsAffected != 0 {
			return fiber.NewError(fiber.StatusBadRequest, `{ "message": "reservation is already done in this date, select other time" }`)
		}

		reservation.From = request.From
		reservation.To = request.To
		reservation.Remarks = request.Remarks
		reservation.YachtID = request.YachtID
		reservation.AccountID = uint(accountID)
		reservation.Created = time.Now()

		if result := create(&reservation); result.Error != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, result.Error.Error())
		}

		response, err := json.Marshal(reservation)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(response)
	}
}
