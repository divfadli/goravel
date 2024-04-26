package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Register struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
	Name     string `json:"name" form:"name"`
	Nik      string `json:"nik" form:"nik"`
	UserType string `json:"user_type" form:"user_type"`
}

func (r *Register) Authorize(ctx http.Context) error {
	return nil
}

func (r *Register) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"username": "required",
		"password": "required|min_len:8",
		"email":    "required|email",
		"name":     "required",
		"nik":      "required",
	}
}

func (r *Register) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"username.required": "Username Cannot be Empty",
		"password.required": "Password Cannot be Empty",
		"password.min_len":  "The password must be at least 8 characters",
		"email.required":    "Email Cannot be Empty",
		"email.email":       "Invalid Email Format",
		"name.required":     "Name Cannot be Empty",
		"nik.required":      "Nik Cannot be Empty",
	}
}

func (r *Register) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Register) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
