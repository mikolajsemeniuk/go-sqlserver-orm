package yachts

import (
	"time"
)

type Account struct {
	Key     string    `json:"key"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Yacht struct {
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Price       float32   `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type Reservation struct {
	Key        string    `json:"key"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Remarks    string    `json:"remarks"`
	AccountKey uint      `json:"accountKey"`
	YachtKey   uint      `json:"yachtKey"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
