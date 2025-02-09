package dto

import "time"

type UserMatchDTO struct {
	ID        int        `json:"ID"`
	UserId    int        `json:"user_id"`
	MatchId   int        `json:"match_id"`
	OrderId   int        `json:"order_id"`
	MediaId   int        `json:"media_id"`
	URL       string     `json:"url"`
	ChatId    *string    `json:"chat_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
