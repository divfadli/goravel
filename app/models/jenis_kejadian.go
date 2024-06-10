package models

import (
	"github.com/goravel/framework/database/orm"
)

type JenisKejadian struct {
	IDJenisKejadian string `json:"id_jenis_kejadian" gorm:"primary_key;column:id_jenis_kejadian"`
	NamaKejadian    string `json:"nama_kejadian" gorm:"not null;column:nama_kejadian"`
	KlasifikasiName string `json:"klasifikasi_name" gorm:"not null;column:klasifikasi_name"`
	orm.Timestamps
	CreatedBy string `json:"created_by" gorm:"not null;column:created_by"`
	orm.SoftDeletes
}

func (r *JenisKejadian) TableName() string {
	return "public.jenis_kejadian"
}
