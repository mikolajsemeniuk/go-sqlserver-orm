package admin

//go:generate goyacc -o gopher.go -p parser gopher.y
type CreateYachtRequest struct {
	Name        string  `json:"name" example:"Maxus"`
	Type        string  `json:"type" example:"yacht"`
	Price       float32 `json:"price" example:"100"`
	Description string  `json:"description" example:"your favorite yacht"`
	Image       string  `json:"image" example:"https://images.pexels.com/photos/843633/pexels-photo-843633.jpeg?auto=compress&cs=tinysrgb&w=1200"`
}
