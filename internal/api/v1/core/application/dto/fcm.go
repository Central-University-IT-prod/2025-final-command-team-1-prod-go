package dto

type NotificationDto struct {
	Title   string `json:"title" binding:"required,max=300,min=1"`
	Content string `json:"content" binding:"required,max=300,min=1"`
	Image   string `json:"image"`
}
