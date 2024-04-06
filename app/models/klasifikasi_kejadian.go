package models

import (
	"github.com/goravel/framework/database/orm"
)

type KlasifikasiKejadian struct {
	IDKlasifikasi uint `gorm:"primaryKey" column:"id_klasifikasi"`
	NamaKlasifikasi string `gorm:"not null" column:"nama_klasifikasi"`
	orm.Timestamps
}

func (r *KlasifikasiKejadian) TableName() string {
	return "rekapitulasi.klasifikasi_kejadian"
}