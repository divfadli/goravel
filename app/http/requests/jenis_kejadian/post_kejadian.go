package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostKejadian struct {
	IdJenisKejadian string `form:"id_jenis_kejadian" json:"id_jenis_kejadian" binding:"required"`
	NamaKejadian    string `form:"nama_kejadian" json:"nama_kejadian" binding:"required"`
	KlasifikasiName string `form:"klasifikasi_name" json:"klasifikasi_name" binding:"required"`
}

func (r *PostKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_name": "required",
	}
}

func (r *PostKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_name.required": "Klasifikasi Name tidak boleh kosong!!",
	}
}

func (r *PostKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
