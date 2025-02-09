package model

import "time"

func (UserMatch) TableName() string {
	return "user_match"
}

type UserMatch struct {
	ID         int        `json:"id" gorm:"id,primaryKey"`
	UserID     int        `json:"user_id" gorm:"user_id"`
	MatchID    int        `json:"match_id" gorm:"match_id"`
	Match_type int        `json:"match_type" gorm:"match_type"`
	ChatID     *string    `json:"chat_id" gorm:"chat_id"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (UserMatchSchema) TableName() string {
	return "user_match"
}

type UserMatchSchema struct {
	ID        int        `json:"id" gorm:"id,primaryKey;autoIncrement"`
	UserId    int        `json:"user_id" gorm:"user_id;index:user_id_match_id;index:match_id_user_id"`
	MatchId   int        `json:"match_id" gorm:"match_id;index:match_id_user_id;index:user_id_match_id"`
	MatchType int        `json:"match_type" gorm:"match_type"`
	ChatID    string     `json:"chat_id" gorm:"chat_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type UserLikers struct {
	UserID int    `json:"user_id"`
	Image  string `json:"image"`
}
