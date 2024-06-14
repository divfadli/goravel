package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostRole struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func (r *PostRole) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostRole) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name": "required",
	}
}

func (r *PostRole) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"name.required": "Name Tidak Boleh Kosong!!",
	}
}

func (r *PostRole) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostRole) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
