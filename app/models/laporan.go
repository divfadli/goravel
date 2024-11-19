package models

import (
	"github.com/goravel/framework/database/orm"
)

type Laporan struct {
	IDLaporan    int64  `json:"id_laporan" gorm:"primary_key" column:"id_laporan"`
	NamaLaporan  string `json:"nama_laporan" gorm:"default:not null" column:"nama_laporan"`
	JenisLaporan string `json:"jenis_laporan" gorm:"default:not null" column:"jenis_laporan"`
	MingguKe     int    `json:"minggu_ke" gorm:"default:0" column:"minggu_ke"`
	BulanKe      int    `json:"bulan_ke" gorm:"default:0" column:"bulan_ke"`
	TahunKe      int    `json:"tahun_ke" gorm:"default:0" column:"tahun_ke"`
	orm.Timestamps
	Dokumen string `json:"dokumen" gorm:"default:not null" column:"dokumen"`
}

func (r *Laporan) TableName() string {
	return "public.laporan"
}
