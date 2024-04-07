package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type RekapKejadianDataKeamanan struct {
	IdRekapKeamanan   uint8       `gorm:"primary_key" column:"id_rekap_keamanan"`
	Tanggal           carbon.Date `gorm:"default:not null" json:"tanggal"`
	TypeKejadianId    string      `gorm:"default:not null" column:"type_kejadian_id"`
	NamaKapal         string      `gorm:"default:not null" column:"nama_kapal"`
	SumberBerita      string      `gorm:"default:not null" column:"sumber_berita"`
	LinkBerita        string      `gorm:"default:not null" column:"link_berita"`
	LokasiKejadian    string      `gorm:"default:not null" column:"lokasi_kejadian"`
	Muatan            string      `gorm:"default:not null" column:"muatan"`
	Asal              *string     `gorm:"default:null" column:"asal"`
	Bendera           *string     `gorm:"default:null" column:"bendera"`
	Tujuan            *string     `gorm:"default:null" column:"tujuan"`
	Latitude          float64     `gorm:"default:0" column:"latitude"`
	Longitude         float64     `gorm:"default:0" column:"longitude"`
	KategoriSumber    string      `gorm:"default:not null" column:"kategori_sumber"`
	TindakLanjut      string      `gorm:"default:not null" column:"tindak_lanjut"`
	IMOKapal          *string     `gorm:"default:null" column:"imo_kapal"`
	MMSIKapal         *string     `gorm:"default:null" column:"mmsi_kapal"`
	InformasiKategori string      `gorm:"default:not null" column:"informasi_kategori"`
	Zona              string      `gorm:"default:not null" column:"zona"`
	CreatedBy         string      `gorm:"default:not null" column:"created_by"`
	IsLocked          bool        `gorm:"default:false" column:"is_locked"`
	orm.Timestamps
}

func (r *RekapKejadianDataKeamanan) TableName() string {
	return "rekapitulasi.rekap_kejadian_data_keamanan"
}
