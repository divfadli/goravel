package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type Karyawan struct {
	EmpNo        string      `json:"emp_no" gorm:"primary_key"`
	Name         string      `json:"name" gorm:"default:not null" column:"name"`
	Gender       string      `json:"gender" gorm:"default:not null" column:"gender"`
	Agama        string      `json:"agama" gorm:"default:not null" column:"agama"`
	TanggalLahir carbon.Date `json:"tanggal_lahir" gorm:"default:not null" column:"tanggal_lahir"`
	UserID       int         `json:"user_id" gorm:"default:0" column:"user_id"`
	User         Akun        `json:"user" gorm:"foreign_key:user_id"`
	JabatanID    int         `json:"jabatan_id" gorm:"default:0" column:"jabatan_id"`
	Jabatan      Jabatan     `json:"jabatan" gorm:"foreign_key:JabatanID"`
	IDAtasan     *string     `json:"id_atasan" gorm:"default: null" column:"id_atasan"`
	orm.Timestamps
}

func (r *Karyawan) TableName() string {
	return "public.karyawan"
}
