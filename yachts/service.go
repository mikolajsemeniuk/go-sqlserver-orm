package yachts

import "time"

type Account struct {
	Key     string    `json:"key" gorm:"primaryKey"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Yacht struct {
	Key         string    `json:"key" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Price       float32   `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type Reservation struct {
	Key        string    `json:"key" gorm:"primaryKey"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Remarks    string    `json:"remarks"`
	AccountKey string    `json:"accountKey"`
	YachtKey   string    `json:"yachtKey"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
