package approval

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostApproval struct {
	Nik        string `form:"nik" json:"nik"`
	Status     string `form:"status" json:"status"`
	IdLaporan  uint8  `form:"id_laporan" json:"id_laporan"`
	Keterangan string `form:"keterangan" json:"keterangan"`
}

func (r *PostApproval) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostApproval) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostApproval) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostApproval) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostApproval) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
