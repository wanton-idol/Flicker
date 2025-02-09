package dto

import "time"

type EventResponseDTO struct {
	ID          uint      `json:"id"`
	EventTime   time.Time `json:"date_time"`
	Location    Location  `json:"location"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Attendees   string    `json:"attendees"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}
