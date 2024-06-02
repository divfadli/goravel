package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type KejadianKeamanan struct {
	IdKejadianKeamanan uint8         `json:"id_kejadian_keamanan" gorm:"primary_key;column:id_kejadian_keamanan"`
	Tanggal            carbon.Date   `json:"tanggal" gorm:"default:not null"`
	JenisKejadianId    string        `json:"jenis_kejadian_id" gorm:"default:not null;column:jenis_kejadian_id"`
	JenisKejadian      JenisKejadian `gorm:"foreign_key:JenisKejadianId;references:IDJenisKejadian"`
	NamaKapal          string        `json:"nama_kapal" gorm:"default:not null;column:nama_kapal"`
	SumberBerita       string        `json:"sumber_berita" gorm:"default:not null;column:sumber_berita"`
	LinkBerita         string        `json:"link_berita" gorm:"default:not null;column:link_berita"`
	LokasiKejadian     string        `json:"lokasi_kejadian" gorm:"default:not null;column:lokasi_kejadian"`
	Muatan             string        `json:"muatan" gorm:"default:not null;column:muatan"`
	Asal               *string       `json:"asal" gorm:"default:null;column:asal"`
	Bendera            *string       `json:"bendera" gorm:"default:null;column:bendera"`
	Tujuan             *string       `json:"tujuan" gorm:"default:null;column:tujuan"`
	Latitude           float64       `json:"latitude" gorm:"default:0;column:latitude"`
	Longitude          float64       `json:"longitude" gorm:"default:0;column:longitude"`
	KategoriSumber     string        `json:"kategori_sumber" gorm:"default:not null;column:kategori_sumber"`
	TindakLanjut       string        `json:"tindak_lanjut" gorm:"default:not null;column:tindak_lanjut"`
	IMOKapal           *string       `json:"imo_kapal" gorm:"default:null;column:imo_kapal"`
	MMSIKapal          *string       `json:"mmsi_kapal" gorm:"default:null;column:mmsi_kapal"`
	InformasiKategori  string        `json:"informasi_kategori" gorm:"default:not null;column:informasi_kategori"`
	Zona               string        `json:"zona" gorm:"default:not null;column:zona"`
	IsLocked           bool          `json:"is_locked" gorm:"default:false;column:is_locked"`
	orm.Timestamps
	CreatedBy string `json:"created_by" gorm:"default:not null;column:created_by"`
}

// type DetailKejadianKeamanan struct {
// 	KejadianKeamanan
// 	JenisKejadian
// }

func (r *KejadianKeamanan) TableName() string {
	return "public.kejadian_keamanan"
}
