package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Name     string `gorm:"default:null" column:"name"`
	Email    string `gorm:"default:null" column:"email"`
	Nik      string `gorm:"unique;not null" column:"nik"`
	Username string `gorm:"not null" column:"username"`
	Password string `gorm:"not null" column:"password"`
	UserType string `gorm:"default:null" column:"user_type"`
}

func (r *User) TableName() string {
	return "public.users"
}
