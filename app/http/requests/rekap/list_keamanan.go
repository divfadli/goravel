package rekap

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ListKeamanan struct {
	Nik          string `form:"nik" json:"nik" binding:"required"`
	Key          string `form:"key" json:"key"`
	TanggalAwal  string `form:"tanggal_awal" json:"tanggal_awal"`
	TanggalAkhir string `form:"tanggal_akhir" json:"tanggal_akhir"`
	Zona         string `form:"zona" json:"zona"`
}

func (r *ListKeamanan) Authorize(ctx http.Context) error {
	return nil
}

func (r *ListKeamanan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nik": "required",
	}
}

func (r *ListKeamanan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nik.required": "NIK Cannot be Empty",
	}
}

func (r *ListKeamanan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListKeamanan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
