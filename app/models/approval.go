package models

import "github.com/golang-module/carbon/v2"

type Approval struct {
	IDApproval         uint8           `json:"id_approval" gorm:"primary_key" column:"id_approval"`
	LaporanID          int64           `json:"laporan_id" gorm:"default: not null" column:"laporan_id"`
	Laporan            Laporan         `gorm:"foreign_key:LaporanID"`
	Status             string          `json:"status" gorm:"default: not null" column:"status"`
	ApprovedBy         string          `json:"approved_by" gorm:"default: not null" column:"approved_by"`
	ApprovedByKaryawan Karyawan        `gorm:"foreign_key:ApprovedBy"`
	Keterangan         *string         `json:"keterangan" gorm:"default: null" column:"keterangan"`
	CreatedAt          carbon.DateTime `gorm:"autoCreateTime;column:created_at"`
}

func (r *Approval) TableName() string {
	return "public.approval"
}
