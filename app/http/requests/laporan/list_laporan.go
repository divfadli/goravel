package laporan

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ListLaporan struct {
	JenisLaporan string `form:"jenis_laporan" json:"jenis_laporan"`
	Minggu       int    `form:"minggu" json:"minggu"`
	Bulan        int    `form:"bulan" json:"bulan"`
	Tahun        int    `form:"tahun" json:"tahun"`
}

func (r *ListLaporan) Authorize(ctx http.Context) error {
	return nil
}

func (r *ListLaporan) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListLaporan) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListLaporan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ListLaporan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
