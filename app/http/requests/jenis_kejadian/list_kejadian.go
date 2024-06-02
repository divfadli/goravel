package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ListKejadian struct {
	KlasifikasiName string `form:"klasifikasi_name" json:"klasifikasi_name" binding:"required"`
}

func (r *ListKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *ListKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_name": "required",
	}
}

func (r *ListKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_name.required": "Klasifikasi Name tidak boleh kosong!!",
	}
}

func (r *ListKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
