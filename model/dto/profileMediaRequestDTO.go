package dto

type ProfileMediaRequestDTO struct {
	ImageText string  `form:"image_text"`
	City      string  `form:"city" `
	Latitude  float64 `form:"latitude"`
	Longitude float64 `form:"longitude"`
	OrderId   int     `form:"order_id"`
}
