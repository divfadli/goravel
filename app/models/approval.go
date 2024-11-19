package models

import "github.com/golang-module/carbon/v2"

type Approval struct {
	IDApproval int64           `json:"id_approval" gorm:"primary_key" column:"id_approval"`
	LaporanID  int64           `json:"laporan_id" gorm:"default: 0" column:"laporan_id"`
	Laporan    Laporan         `gorm:"foreign_key:LaporanID"`
	Status     string          `json:"status" gorm:"default: not null" column:"status"`
	ApprovedBy string          `json:"approved_by" gorm:"default: not null" column:"approved_by"`
	Karyawan   Karyawan        `gorm:"foreignKey:ApprovedBy;references:EmpNo"`
	Keterangan *string         `json:"keterangan" gorm:"default: null" column:"keterangan"`
	CreatedAt  carbon.DateTime `gorm:"autoCreateTime;column:created_at"`
}

func (r *Approval) TableName() string {
	return "public.approval"
}
