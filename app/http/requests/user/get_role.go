package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetRole struct {
	Name string `form:"id_user" json:"id_user" binding:"required"`
}

func (r *GetRole) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetRole) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"id_user": "required",
	}
}

func (r *GetRole) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"id_user.required": "ID User Tidak Boleh Kosong!!",
	}
}

func (r *GetRole) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetRole) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
