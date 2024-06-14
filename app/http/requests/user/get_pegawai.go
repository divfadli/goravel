package user

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetPegawai struct {
	KeyPegawai string `form:"key_pegawai" json:"key_pegawai" binding:"required"`
}

func (r *GetPegawai) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetPegawai) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"key_pegawai":           "required",
	}
}

func (r *GetPegawai) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"key_pegawai.required":           "Kata Kunci Tidak Boleh Kosong!!",
	}
}

func (r *GetPegawai) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetPegawai) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
