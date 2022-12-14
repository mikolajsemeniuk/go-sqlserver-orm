package admin

import "time"

type Yacht struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Price       float32    `json:"price"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Updated     *time.Time `json:"updated"`
}
