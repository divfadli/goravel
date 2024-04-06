package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type RekapKejadianKeamananController struct {
	//Dependent services
}

func NewRekapKejadianKeamananController() *RekapKejadianKeamananController {
	return &RekapKejadianKeamananController{
		//Inject services
	}
}

func (r *RekapKejadianKeamananController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *RekapKejadianKeamananController) ListRekapKeamanan(ctx http.Context) http.Response {
	// var req rekap.ListKeamanan
	// if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
	// 	return SanitizeGet(ctx, chekRequestErr)
	// }

	// var rekap_keamanan models.RekapKejadianDataKeamanan
	// query := facades.Orm().Query().Table("rekapitulasi.rekap_kejadian_data_keamanan").
	// 	Join(`join rekapitulasi.kejadian ON rekapitulasi.kejadian.id_type_kejadian = rekapitulasi.rekap_kejadian_data_keamanan.type_kejadian_id
	// 		  join rekapitulasi.klasifikasi_kejadian ON rekapitulasi.kejadian.klasifikasi_id = rekapitulasi.klasifikasi_kejadian.id_klasifikasi`)

	// if req.Key != "" {
	// 	query = query.Where("rekapitulasi.rekap_kejadian_data_keamanan.keamanan_id", req.Key)
	// }
	// if err := facades.Orm().Query().Table("rekapitulasi.kejadian").
	// 	Join("join rekapitulasi.klasifikasi_kejadian ON rekapitulasi.kejadian.klasifikasi_id = rekapitulasi.klasifikasi_kejadian.id_klasifikasi").
	// 	Where("id_type_kejadian", req.IdTypeKejadian).
	// 	Scan(&data); err != nil {
	// 	return ErrorSystem(ctx, "Data Tidak Ada")
	// }

	return nil
}
