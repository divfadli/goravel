package models

import (
	"github.com/goravel/framework/database/orm"
)

type Akun struct {
	IDUser   int64  `json:"id_user" gorm:"primary_key" column:"id_user"`
	Email    string `gorm:"default:not null" column:"email"`
	Password string `gorm:"default:not null" column:"password"`
	orm.Timestamps
	RoleID int  `json:"role_id" gorm:"default:null" column:"role_id"`
	Role   Role `gorm:"foreign_key:RoleID"`
	orm.SoftDeletes
}

func (r *Akun) TableName() string {
	return "public.akun"
}
