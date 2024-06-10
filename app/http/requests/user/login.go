package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Login struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (r *Login) Authorize(ctx http.Context) error {
	return nil
}

func (r *Login) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required|min_len:8",
	}
}

func (r *Login) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":    "Email Tidak Boleh Kosong ",
		"email.email":       "Email Tidak Valid",
		"password.required": "Password Tidak Boleh Kosong ",
		"password.min_len":  "Masukkan Password Minimal 8 Karakter",
	}
}

func (r *Login) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Login) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
