package models

import (
	"github.com/goravel/framework/database/orm"
)

type LaporanKeamanan struct {
	IDLaporanKeamanan  int64            `json:"id_laporan_keamanan" gorm:"primary_key" column:"id_laporan_keamanan"`
	LaporanID          int              `json:"laporan_id" gorm:"default: not null" column:"laporan_id"`
	Laporan            Laporan          `gorm:"foreign_key:LaporanID"`
	KejadianKeamananID int              `json:"kejadian_keamanan_id" gorm:"default: not null" column:"kejadian_keamanan_id"`
	KejadianKeamanan   KejadianKeamanan `gorm:"foreign_key:KejadianKeamananID"`
	orm.Timestamps
}

func (r *LaporanKeamanan) TableName() string {
	return "public.laporan_keamanan"
}
