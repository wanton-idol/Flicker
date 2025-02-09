package model

import "gorm.io/gorm"

type ChatDetails struct {
	gorm.Model
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	MediaURL   string `json:"media_url"`
	ChatID     string `json:"chat_id"`
	IsRead     bool   `json:"is_read"`
}

type ChatValues struct {
	ChatID      string      `json:"chat_id"`
	UserProfile UserProfile `json:"user_profile"`
	Image       string      `json:"image"`
}

type MessagesIDs struct {
	MessageIDS []int `json:"message_ids"`
}

type ChatIDs struct {
	ChatIDS []string `json:"chat_ids"`
}
