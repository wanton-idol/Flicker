package model

import "gorm.io/gorm"

type Nudge struct {
	ID       int    `json:"id" gorm:"column:ID"`
	Question string `json:"question" gorm:"column:question"`
}

type UserNudge struct {
	gorm.Model
	UserID   int    `json:"user_id" gorm:"column:user_id"`
	Question string `json:"question" gorm:"column:question"`
	Answer   string `json:"answer" gorm:"column:answer"`
	Order    int    `json:"order" gorm:"column:order"`
	MediaURL string `json:"media_url" gorm:"column:media_url"`
	Type     string `json:"type" gorm:"column:type"`
}

type NudgeDetail struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Order    int    `json:"order"`
	MediaURL string `json:"media_url"`
	Type     string `json:"type"`
}

type NudgeRequest struct {
	Question string `form:"question"`
	Answer   string `form:"answer"`
	Order    int    `form:"order"`
	Type     string `form:"type"`
}
