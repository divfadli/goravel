package models

import (
	"github.com/goravel/framework/database/orm"
)

type Role struct {
	RoleID int    `gorm:"primary_key" json:"role_id"`
	Name   string `gorm:"not null" json:"name"`
	EmpNo  string `gorm:"not null" json:"emp_no"`
	orm.Timestamps
}

type MyRole struct {
	Name string `json:"name"`
}

func (r *Role) TableName() string {
	return "public.role"
}
