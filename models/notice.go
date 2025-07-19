package models

import "time"

type Notice struct {
	NoticeID    uint64    `gorm:"primaryKey" json:"noticeID"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Username    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	IsImportant bool      `json:"is_important"`
	Role        string    `json:"role"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Notice) TableName() string {
	return "notices"
}
