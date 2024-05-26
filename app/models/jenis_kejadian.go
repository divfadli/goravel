package models

import (
	"github.com/goravel/framework/database/orm"
)

type JenisKejadian struct {
	IDJenisKejadian string `json:"id_jenis_kejadian" gorm:"primary_key" column:"id_jenis_kejadian"`
	NamaKejadian    string `json:"nama_kejadian" gorm:"default:not null" column:"nama_kejadian"`
	KlasifikasiName string `json:"klasifikasi_name" gorm:"default:not null" column:"klasifikasi_name"`
	orm.Timestamps
}

func (r *JenisKejadian) TableName() string {
	return "public.jenis_kejadian"
}
