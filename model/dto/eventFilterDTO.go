package dto

import "time"

type EventFilterDTO struct {
	Type      string     `json:"type"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Distance  int        `json:"distance"`
}
