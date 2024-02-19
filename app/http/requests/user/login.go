package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type Login struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (r *Login) Authorize(ctx http.Context) error {
	return nil
}

func (r *Login) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"username": "required",
		"password": "required|min_len:8",
	}
}

func (r *Login) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"username.required": "Username Cannot be Empty",
		"password.required": "Password Cannot be Empty",
		"password.min_len":  "The password must be at least 8 characters",
	}
}

func (r *Login) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *Login) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
