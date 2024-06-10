package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostKejadian struct {
	IdJenisKejadian string `form:"id_jenis_kejadian" json:"id_jenis_kejadian"`
	NamaKejadian    string `form:"nama_kejadian" json:"nama_kejadian" binding:"required"`
	KlasifikasiName string `form:"klasifikasi_name" json:"klasifikasi_name" binding:"required"`
	Nik             string `form:"nik" json:"nik" binding:"required"`
}

func (r *PostKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nama_kejadian":    "required",
		"klasifikasi_name": "required",
		"nik":              "required",
	}
}

func (r *PostKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nama_kejadian.required":    "Nama Kejadian Tidak Boleh Kosong!!",
		"klasifikasi_name.required": "Klasifikasi Nama Tidak Boleh Kosong!!",
		"nik.required":              "NIK Tidak Boleh Kosong!!",
	}
}

func (r *PostKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
