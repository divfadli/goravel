package rekap

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetKeamanan struct {
	IdRekapKeamanan string `form:"id_rekap_keamanan" json:"id_rekap_keamanan" binding:"required"`
}

func (r *GetKeamanan) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetKeamanan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_rekap_keamanan": "required",
	}
}

func (r *GetKeamanan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_rekap_keamanan.required": "ID Rekap Keamanan Cannot be Empty",
	}
}

func (r *GetKeamanan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetKeamanan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
