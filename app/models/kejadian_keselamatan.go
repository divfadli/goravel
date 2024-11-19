package models

import (
	"encoding/json"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type KejadianKeselamatan struct {
	IdKejadianKeselamatan int64           `json:"id_kejadian_keselamatan" gorm:"primary_key" column:"id_kejadian_keselamatan"`
	Tanggal               carbon.Date     `json:"tanggal" gorm:"default:not null"`
	JenisKejadianId       string          `json:"jenis_kejadian_id" gorm:"default:not null" column:"jenis_kejadian_id"`
	JenisKejadian         JenisKejadian   `json:"jenis_kejadian" gorm:"foreign_key:JenisKejadianId;references:IDJenisKejadian"`
	NamaKapal             string          `json:"nama_kapal" gorm:"default:not null" column:"nama_kapal"`
	SumberBerita          string          `json:"sumber_berita" gorm:"default:not null" column:"sumber_berita"`
	LinkBerita            string          `json:"link_berita" gorm:"default:not null" column:"link_berita"`
	LokasiKejadian        string          `json:"lokasi_kejadian" gorm:"default:not null" column:"lokasi_kejadian"`
	Korban                json.RawMessage `json:"korban" gorm:"default:not null"  column:"korban"`
	Latitude              float64         `json:"latitude" gorm:"default:0" column:"latitude"`
	Longitude             float64         `json:"longitude" gorm:"default:0" column:"longitude"`
	Penyebab              string          `json:"penyebab" gorm:"default:not null" column:"penyebab"`
	TipeSumberKejadian    string          `json:"tipe_sumber_kejadian" gorm:"default:not null" column:"tipe_sumber_kejadian"`
	PelabuhanAsal         string          `json:"pelabuhan_asal" gorm:"default:null" column:"pelabuhan_asal"`
	PelabuhanTujuan       string          `json:"pelabuhan_tujuan" gorm:"default:null" column:"pelabuhan_tujuan"`
	TindakLanjut          string          `json:"tindak_lanjut" gorm:"default:not null" column:"tindak_lanjut"`
	Keterangan            string          `json:"keterangan" gorm:"default:not null" column:"keterangan"`
	Zona                  string          `json:"zona" gorm:"default:not null" column:"zona"`
	IsLocked              bool            `json:"is_locked" gorm:"default:false" column:"is_locked"`
	orm.Timestamps
	CreatedBy string `json:"created_by" gorm:"default: not null" column:"created_by"`
}

type KejadianKeselamatanImage struct {
	KejadianKeselamatan
	FileImage []FileImage `json:"file_image"`
}

type KejadianKeselamatanKorban struct {
	KejadianKeselamatan `json:"kejadian_keselamatan"`
	ListKorban          `json:"list_korban"`
}

type ListKorban struct {
	KorbanTewas   int `json:"korban_tewas"`
	KorbanSelamat int `json:"korban_selamat"`
	KorbanHilang  int `json:"korban_hilang"`
}

func (r *KejadianKeselamatan) TableName() string {
	return "public.kejadian_keselamatan"
}
