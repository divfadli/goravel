package models

import (
	"github.com/goravel/framework/database/orm"
)

type Kejadian struct {
	IDTypeKejadian   string `gorm:"unique;primary_key" column:"id_type_kejadian"`
	JenisPelanggaran string `gorm:"not null;unique" column:"jenis_pelanggaran"`
	KlasifikasiID    uint   `gorm:"not null" column:"klasifikasi_id"`
	orm.Timestamps
}

func (r *Kejadian) TableName() string {
	return "rekapitulasi.kejadian"
}
