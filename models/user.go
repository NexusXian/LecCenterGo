package models

import (
	"LecCenterGo/dao"
	_ "gorm.io/gorm"
	"time"
)
//define User struct
type User struct {
	UserID    uint64    `gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Grade     string    `json:"grade"`
	Major     string    `json:"major"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role" gorm:"default:'user'"`
	Count     int       `json:"count"`
	Direction string    `json:"direction"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

//Get user by Email from database
func GetUserByEmail(email string) (*User, error) {
	var user User
	result := dao.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByPhone(phone string) (*User, error) {
	var user User
	result := dao.DB.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmailOrPhone(account string) (*User, error) {
	var user User
	result := dao.DB.Where("phone = ? or email = ?", account, account).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(user *User) error {
	result := dao.DB.Create(user)
	return result.Error
}

func FindAllUsers() ([]User, error) {
	var users []User
	result := dao.DB.Find(&users)
	return users, result.Error
}
