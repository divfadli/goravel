package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetKejadian struct {
	IdJenisKejadian string `form:"id_jenis_kejadian" json:"id_jenis_kejadian" binding:"required"`
}

func (r *GetKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_jenis_kejadian": "required",
	}
}

func (r *GetKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_jenis_kejadian.required": "ID Jenis Kejadian Tidak Boleh Kosong!!",
	}
}

func (r *GetKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
