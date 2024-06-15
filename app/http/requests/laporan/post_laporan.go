package laporan

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostLaporan struct {
	NamaLaporan  string `form:"nama_laporan" json:"nama_laporan"`
	JenisLaporan string `form:"jenis_laporan" json:"jenis_laporan"`
	MingguKe     int    `form:"minggu_ke" json:"minggu_ke"`
	BulanKe      int    `form:"bulan_ke" json:"bulan_ke"`
	TahunKe      int    `form:"tahun_ke" json:"tahun_ke"`
	CreatedBy    string `form:"created_by" json:"created_by"`
}

func (r *PostLaporan) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostLaporan) Rules(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostLaporan) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostLaporan) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostLaporan) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
