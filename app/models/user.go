package models

import "github.com/golang-module/carbon/v2"

type User struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"default:null" json:"name"`
	Email         string          `gorm:"default:null" json:"email"`
	Nik           string          `gorm:"unique;not null" json:"nik"`
	Username      string          `gorm:"not null" json:"username"`
	Password      string          `gorm:"not null" json:"password"`
	RememberToken string          `gorm:"default:null" json:"remember_token"`
	Type          string          `gorm:"default:null" json:"type"`
	UserType      string          `gorm:"default:null" json:"user_type"`
	CreatedAt     carbon.DateTime `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt     carbon.DateTime `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "public.users"
}
