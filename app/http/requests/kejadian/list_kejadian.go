package kejadian

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ListKejadian struct {
	KlasifikasiId int `form:"klasifikasi_id" json:"klasifikasi_id" binding:"required"`
}

func (r *ListKejadian) Authorize(ctx http.Context) error {
	return nil
}

func (r *ListKejadian) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_id": "required",
	}
}

func (r *ListKejadian) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"klasifikasi_id.required": "Klasifikasi ID Cannot be Empty",
	}
}

func (r *ListKejadian) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListKejadian) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
