package models

import "time"

type Attendance struct {
	ID        uint      `gorm:"primaryKey" json:"ID"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Img       string    `json:"img"`
	CheckInfo string    `json:"check_info"`
}

func (Attendance) TableName() string {
	return "attendances"
}
