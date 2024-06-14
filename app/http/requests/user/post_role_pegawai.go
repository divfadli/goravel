package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostRolePegawai struct {
	IDUser uint8  `form:"id_user" json:"id_user"`
	EmpNo  string `form:"emp_no" json:"emp_no" binding:"required"`
	RoleID int    `form:"role_id" json:"role_id" binding:"required"`
}

func (r *PostRolePegawai) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostRolePegawai) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"emp_no":  "required",
		"role_id": "required",
	}
}

func (r *PostRolePegawai) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"emp_no.required":  "ID User Tidak Boleh Kosong!!",
		"role_id.required": "Role ID Tidak Boleh Kosong!!",
	}
}

func (r *PostRolePegawai) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostRolePegawai) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
