package kejadian_keselamatan

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostKeselamatan struct {
	IdKejadianKeselamatan int     `form:"id_kejadian_keselamatan" json:"id_kejadian_keselamatan"`
	Tanggal               string  `form:"tanggal" json:"tanggal" binding:"required"`
	JenisKejadianId       string  `form:"jenis_kejadian_id" json:"jenis_kejadian_id" binding:"required"`
	NamaKapal             string  `form:"nama_kapal" json:"nama_kapal"`
	SumberBerita          string  `form:"sumber_berita" json:"sumber_berita" binding:"required"`
	LinkBerita            string  `form:"link_berita" json:"link_berita" binding:"required"`
	LokasiKejadian        string  `form:"lokasi_kejadian" json:"lokasi_kejadian" binding:"required"`
	KorbanTewas           int     `form:"korban_tewas" json:"korban_tewas"`
	KorbanSelamat         int     `form:"korban_selamat" json:"korban_selamat"`
	KorbanHilang          int     `form:"korban_hilang" json:"korban_hilang"`
	Latitude              float64 `form:"latitude" json:"latitude" binding:"required"`
	Longitude             float64 `form:"longitude" json:"longitude" binding:"required"`
	Penyebab              string  `form:"penyebab" json:"penyebab" binding:"required"`
	PelabuhanAsal         string  `form:"pelabuhan_asal" json:"pelabuhan_asal"`
	PelabuhanTujuan       string  `form:"pelabuhan_tujuan" json:"pelabuhan_tujuan"`
	KategoriSumber        string  `form:"kategori_sumber" json:"kategori_sumber" binding:"required"`
	TindakLanjut          string  `form:"tindak_lanjut" json:"tindak_lanjut" binding:"required"`
	Keterangan            string  `form:"keterangan" json:"keterangan" binding:"required"`
	Zona                  string  `form:"zona" json:"zona" binding:"required"`
	Nik                   string  `form:"nik" json:"nik" binding:"required"`
}

func (r *PostKeselamatan) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostKeselamatan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"tanggal":           "required",
		"jenis_kejadian_id": "required",
		"sumber_berita":     "required",
		"link_berita":       "required",
		"lokasi_kejadian":   "required",
		"penyebab":          "required",
		"latitude":          "required",
		"longitude":         "required",
		"kategori_sumber":   "required",
		"tindak_lanjut":     "required",
		"keterangan":        "required",
		"zona":              "required",
		"nik":               "required",
	}
}

func (r *PostKeselamatan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"tanggal.required":           "Tanggal Tidak Boleh Kosong!!",
		"jenis_kejadian_id.required": "Jenis Kejadian ID Tidak Boleh Kosong!!",
		"sumber_berita.required":     "Sumber Berita Tidak Boleh Kosong!!",
		"link_berita.required":       "Link Berita Tidak Boleh Kosong!!",
		"lokasi_kejadian.required":   "Lokasi Kejadian Tidak Boleh Kosong!!",
		"penyebab.required":          "Penyebab Tidak Boleh Kosong!!",
		"latitude.required":          "Latitude Tidak Boleh Kosong!!",
		"longitude.required":         "Longitude Tidak Boleh Kosong!!",
		"kategori_sumber.required":   "Kategori Sumber Tidak Boleh Kosong!!",
		"tindak_lanjut.required":     "Tindak Lanjut Tidak Boleh Kosong!!",
		"keterangan.required":        "Keterangan Tidak Boleh Kosong!!",
		"zona.required":              "Zona Tidak Boleh Kosong!!",
		"nik.required":               "NIK Tidak Boleh Kosong!!",
	}
}

func (r *PostKeselamatan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostKeselamatan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
