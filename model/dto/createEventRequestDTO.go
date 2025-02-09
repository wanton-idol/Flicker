package dto

import "time"

type CreateEventDTO struct {
	ID          uint      `json:"id"`
	EventTime   time.Time `json:"date_time"`
	Location    Location  `json:"location"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Attendees   string    `json:"attendee"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type Location struct {
	Address1  string  `json:"address1"`
	Address2  string  `json:"address2"`
	City      string  `json:"city"`
	Pincode   string  `json:"pincode"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
