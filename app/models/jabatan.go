package models

import (
	"github.com/goravel/framework/database/orm"
)

type Jabatan struct {
	IdJabatan uint8  `json:"id_jabatan" gorm:"primary_key"`
	Name      string `json:"name" gorm:"default:not null" column:"name"`
	orm.Timestamps
}

func (r *Jabatan) TableName() string {
	return "public.jabatan"
}
