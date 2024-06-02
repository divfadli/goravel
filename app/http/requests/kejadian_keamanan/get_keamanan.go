package rekap

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetKeamanan struct {
	IdKejadianKeamanan string `form:"id_kejadian_keamanan" json:"id_kejadian_keamanan" binding:"required"`
}

func (r *GetKeamanan) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetKeamanan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_kejadian_keamanan": "required",
	}
}

func (r *GetKeamanan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_kejadian_keamanan.required": "ID Kejadian Keamanan tidak boleh kosong!!",
	}
}

func (r *GetKeamanan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetKeamanan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
