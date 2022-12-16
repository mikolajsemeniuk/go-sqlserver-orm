package reservations

import "time"

type Request struct {
	From    time.Time `json:"from" example:"2022-12-16T21:00:00.511Z"`
	To      time.Time `json:"to" example:"2022-12-17T21:00:00.511Z"`
	Remarks string    `json:"remarks" example:"please clean it"`
	YachtID uint      `json:"yachtID" example:"4"`
}
