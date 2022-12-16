package reservations

import "time"

type Reservation struct {
	ID        uint       `gorm:"primaryKey"`
	From      time.Time  `json:"from"`
	To        time.Time  `json:"to"`
	Remarks   string     `json:"remarks"`
	AccountID uint       `json:"account_id" gorm:"column:account_id"`
	YachtID   uint       `json:"yacht_id" gorm:"column:yacht_id"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated"`
}
