package model

type GCMNotification struct {
	GCM string `json:"GCM"`
}

type Notification struct {
	Notification NotificationData `json:"notification"`
}

type NotificationData struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}
