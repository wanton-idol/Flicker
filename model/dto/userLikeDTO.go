package dto

type UserLikeDTO struct {
	LikeeID int `json:"likeeID"`
	LikerID int `json:"likerID"`
	Type    int `json:"type"`
}
