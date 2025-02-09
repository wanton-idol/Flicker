package elasticsearchPkg

import (
	"time"
)

type Event struct {
	ID          uint       `json:"id"`
	UserId      int        `json:"user_id"`
	EventTime   time.Time  `json:"event_time"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Attendees   string     `json:"attendees"`
	ExpiresAt   time.Time  `json:"expires_at"`
	Address1    string     `json:"address1"`
	Address2    string     `json:"address2"`
	City        string     `json:"city"`
	State       string     `json:"state"`
	Pincode     string     `json:"pincode"`
	Location    []float32  `json:"location"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
