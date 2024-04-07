package rekap

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostKeamanan struct {
	IdRekapKeamanan   int     `form:"id_rekap_keamanan" json:"id_rekap_keamanan"`
	Tanggal           string  `form:"tanggal" json:"tanggal" binding:"required"`
	TypeKejadianId    string  `form:"type_kejadian_id" json:"type_kejadian_id" binding:"required"`
	NamaKapal         string  `form:"nama_kapal" json:"nama_kapal"`
	SumberBerita      string  `form:"sumber_berita" json:"sumber_berita" binding:"required"`
	LinkBerita        string  `form:"link_berita" json:"link_berita" binding:"required"`
	LokasiKejadian    string  `form:"lokasi_kejadian" json:"lokasi_kejadian" binding:"required"`
	Muatan            string  `form:"muatan" json:"muatan" binding:"required"`
	Asal              string  `form:"asal" json:"asal"`
	Bendera           string  `form:"bendera" json:"bendera"`
	Tujuan            string  `form:"tujuan" json:"tujuan"`
	Latitude          float64 `form:"latitude" json:"latitude" binding:"required"`
	Longitude         float64 `form:"longitude" json:"longitude" binding:"required"`
	KategoriSumber    string  `form:"kategori_sumber" json:"kategori_sumber" binding:"required"`
	TindakLanjut      string  `form:"tindak_lanjut" json:"tindak_lanjut" binding:"required"`
	IMOKapal          string  `form:"imo_kapal" json:"imo_kapal"`
	MMSIKapal         string  `form:"mmsi_kapal" json:"mmsi_kapal"`
	InformasiKategori string  `form:"informasi_kategori" json:"informasi_kategori" binding:"required"`
	Zona              string  `form:"zona" json:"zona" binding:"required"`
	Nik               string  `form:"nik" json:"nik" binding:"required"`
}

func (r *PostKeamanan) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostKeamanan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"tanggal":            "required",
		"type_kejadian_id":   "required",
		"sumber_berita":      "required",
		"link_berita":        "required",
		"lokasi_kejadian":    "required",
		"muatan":             "required",
		"latitude":           "required",
		"longitude":          "required",
		"kategori_sumber":    "required",
		"tindak_lanjut":      "required",
		"informasi_kategori": "required",
		"zona":               "required",
		"nik":                "required",
	}
}

func (r *PostKeamanan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"tanggal.required":            "Tanggal Cannot be Empty",
		"type_kejadian_id.required":   "Type Kejadian ID Cannot be Empty",
		"sumber_berita.required":      "Sumber Berita Cannot be Empty",
		"link_berita.required":        "Link Berita Cannot be Empty",
		"lokasi_kejadian.required":    "Lokasi Kejadian Cannot be Empty",
		"muatan.required":             "Muatan Cannot be Empty",
		"latitude.required":           "Latitude Cannot be Empty",
		"longitude.required":          "Longitude Cannot be Empty",
		"kategori_sumber.required":    "Kategori Sumber Cannot be Empty",
		"tindak_lanjut.required":      "Tindak Lanjut Cannot be Empty",
		"informasi_kategori.required": "Informasi Kategori Cannot be Empty",
		"zona.required":               "Zona Cannot be Empty",
		"nik.required":                "NIK Cannot be Empty",
	}
}

func (r *PostKeamanan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostKeamanan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
