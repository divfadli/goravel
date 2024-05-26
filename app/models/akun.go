package models

import (
	"github.com/goravel/framework/database/orm"
)

type Akun struct {
	IDUser   uint8  `json:"id_user" gorm:"primary_key" column:"id_user"`
	Email    string `gorm:"default:not null" column:"email"`
	Password string `gorm:"default:not null" column:"password"`
	orm.Timestamps
	RoleID int  `json:"role_id" gorm:"default:not null" column:"role_id"`
	Role   Role `gorm:"foreign_key:RoleID"`
}

func (r *Akun) TableName() string {
	return "public.akun"
}
