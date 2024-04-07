package models

import (
	"encoding/json"

	"github.com/golang-module/carbon/v2"
	"github.com/goravel/framework/database/orm"
)

type RekapKejadianDataKeselamatan struct {
	IdRekapKeselamatan uint8           `gorm:"primary_key" column:"id_rekap_keselamatan"`
	Tanggal            carbon.Date     `gorm:"default:not null" column:"tanggal"`
	TypeIDKejadian     string          `gorm:"default:not null" column:"type_id_kejadian"`
	NamaKapal          string          `gorm:"default:not null" column:"nama_kapal"`
	SumberBerita       string          `gorm:"default:not null" column:"sumber_berita"`
	LokasiKejadian     string          `gorm:"default:not null" column:"lokasi_kejadian"`
	Korban             json.RawMessage `gorm:"default:not null  column:"korban"`
	Latitude           float64         `gorm:"default:not null" column:"latitude"`
	Longitude          float64         `gorm:"default:not null" column:"longitude"`
	Penyebab           string          `gorm:"default:not null" column:"penyebab"`
	TipeSumberKejadian string          `gorm:"default:not null" column:"tipe_sumber_kejadian"`
	PelabuhanAsal      string          `gorm:"default:null" column:"pelabuhan_asal"`
	PelabuhanTujuan    string          `gorm:"default:null" column:"pelabuhan_tujuan"`
	TindakLanjut       string          `gorm:"default:not null" column:"tindak_lanjut"`
	Keterangan         string          `gorm:"default:not null" column:"keterangan"`
	Zona               string          `gorm:"default:not null" column:"zona"`
	IsLocked           bool            `gorm:"default:false" column:"is_locked"`
	orm.Timestamps
}

type Korban struct {
	KorbanTewas   int `gorm:"default:0" json:"korban_tewas"`
	KorbanSelamat int `gorm:"default:0" json:"korban_selamat"`
	KorbanHilang  int `gorm:"default:0" json:"korban_hilang"`
}

func (r *RekapKejadianDataKeselamatan) TableName() string {
	return "rekapitulasi.rekap_kejadian_keselamatan"
}

// var atasan []cuti.AtasanApproved
// 	atasan = append(atasan, temp)
// 	// Marshal atasan into JSON
// 	jsonDataAtasan, _ := json.Marshal(atasan)
// 	sck.ApprovedBy = json.RawMessage(jsonDataAtasan)
