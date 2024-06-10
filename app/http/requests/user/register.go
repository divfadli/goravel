package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Register struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (r *Register) Authorize(ctx http.Context) error {
	return nil
}

func (r *Register) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required|min_len:8",
	}
}

func (r *Register) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"email.required":    "Email Cannot be Empty",
		"email.email":       "Invalid Email Format",
		"password.required": "Password Cannot be Empty",
		"password.min_len":  "The password must be at least 8 characters",
	}
}

func (r *Register) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Register) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
