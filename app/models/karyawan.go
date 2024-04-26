package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type Karyawan struct {
	EmpNo        string      `json:"emp_no" gorm:"primary_key"`
	Name         string      `json:"name"`
	Gender       string      `json:"gender"`
	Agama        string      `json:"agama"`
	TanggalLahir carbon.Date `json:"tanggal_lahir"`
	HP           string      `json:"hp"`
	PosID        string      `json:"pos_id"`
	SupPosID     string      `json:"sup_pos_id"`
	SupEmpNo     string      `json:"sup_emp_no"`
	orm.Timestamps
}

func (r *Karyawan) TableName() string {
	return "public.karyawan"
}
