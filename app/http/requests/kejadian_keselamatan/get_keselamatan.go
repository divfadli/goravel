package kejadian_keselamatan

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetKeselamatan struct {
	IdKejadianKeselamatan string `form:"id_kejadian_keselamatan" json:"id_kejadian_keselamatan" binding:"required"`
}

func (r *GetKeselamatan) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetKeselamatan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_kejadian_keselamatan": "required",
	}
}

func (r *GetKeselamatan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_kejadian_keselamatan.required": "ID Kejadian Keselamatan tidak boleh kosong!!",
	}
}

func (r *GetKeselamatan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetKeselamatan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
