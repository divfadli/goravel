package rekap

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ListKeamanan struct {
	Nik          string `form:"nik" json:"nik" binding:"required"`
	Key          string `form:"key" json:"key"`
	TanggalAwal  string `form:"tanggal_awal" json:"tanggal_awal" binding:"required"`
	TanggalAkhir string `form:"tanggal_akhir" json:"tanggal_akhir" binding:"required"`
	Zona         string `form:"zona" json:"zona"`
}

func (r *ListKeamanan) Authorize(ctx http.Context) error {
	return nil
}

func (r *ListKeamanan) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"nik":           "required",
		"tanggal_awal":  "required",
		"tanggal_akhir": "required",
	}
}

func (r *ListKeamanan) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"nik.required":           "NIK Tidak Boleh Kosong!!",
		"tanggal_awal.required":  "Tanggal Awal Tidak Boleh Kosong!!",
		"tanggal_akhir.required": "Tanggal Akhir Tidak Boleh Kosong!!",
	}
}

func (r *ListKeamanan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListKeamanan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
