package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role struct {
	IDRole int    `gorm:"primary_key" json:"id_role" column:"id_role"`
	Name   string `gorm:"not null" json:"name" column:"name"`
	orm.Timestamps
}

// type MyRole struct {
// 	Name string `json:"name"`
// }

func (r *Role) TableName() string {
	return "public.role"
}
