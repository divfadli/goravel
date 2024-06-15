package laporan

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type GetLaporan struct {
	IdLaporan uint8 `form:"id_laporan" json:"id_laporan"`
}

func (r *GetLaporan) Authorize(ctx http.Context) error {
	return nil
}

func (r *GetLaporan) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetLaporan) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetLaporan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *GetLaporan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
