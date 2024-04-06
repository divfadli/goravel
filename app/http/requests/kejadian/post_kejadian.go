package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostKejadian struct {
	IdTypeKejadian string `form:"id_type_kejadian" json:"id_type_kejadian"`
	JenisPelanggaran string `form:"jenis_pelanggaran" json:"jenis_pelanggaran" binding:"required"`
	KlasifikasiId int `form:"klasifikasi_id" json:"klasifikasi_id" binding:"required"`
}

func (r *PostKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_id": "required",
	}
}

func (r *PostKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_id.required": "Klasifikasi ID Cannot be Empty",
	}
}

func (r *PostKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
