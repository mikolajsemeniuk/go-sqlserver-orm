package sailing

import "time"

type Account struct {
	ID       uint       `gorm:"primaryKey"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Password []byte     `json:"password"`
	Role     string     `json:"role"`
	Created  time.Time  `json:"created"`
	Updated  *time.Time `json:"updated"`
}

type Yacht struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Price       float32   `json:"price"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type Reservation struct {
	ID        uint      `gorm:"primaryKey"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
	Remarks   string    `json:"remarks"`
	AccountID uint      `json:"accountID"`
	YachtID   uint      `json:"yachtID"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
