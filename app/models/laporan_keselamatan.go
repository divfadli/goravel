package models

import (
	"github.com/goravel/framework/database/orm"
)

type LaporanKeselamatan struct {
	IDLaporanKeselamatan  int64               `json:"id_laporan_keselamatan" gorm:"primary_key" column:"id_laporan_keselamatan"`
	LaporanID             int                 `json:"laporan_id" gorm:"default: not null" column:"laporan_id"`
	Laporan               Laporan             `gorm:"foreign_key:LaporanID"`
	KejadianKeselamatanID int                 `json:"kejadian_keselamatan_id" gorm:"default: not null" column:"kejadian_keselamatan_id"`
	KejadianKeselamatan   KejadianKeselamatan `gorm:"foreign_key:KejadianKeselamatanID"`
	orm.Timestamps
}

func (r *LaporanKeselamatan) TableName() string {
	return "public.laporan_keselamatan"
}
