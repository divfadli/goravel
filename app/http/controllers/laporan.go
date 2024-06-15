package controllers

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	lp "goravel/app/http/requests/laporan"
	"goravel/app/models"
)

type Laporan struct {
	// Dependent services
}

func NewLaporan() *Laporan {
	return &Laporan{
		// Inject services
	}
}

func (r *Laporan) Index(ctx http.Context) http.Response {
	return nil
}

func (r *Laporan) Show(ctx http.Context) http.Response {
	return nil
}

func (r *Laporan) Create(ctx http.Context) http.Response {
	var req lp.PostLaporan

	if err := ctx.Request().Bind(&req); err != nil {
		return SanitizeGet(ctx, err)
	}

	file, err := ctx.Request().File("files")
	if err != nil {
		return Error(ctx, http.StatusInternalServerError, err.Error())
	}

	path, nameLaporan, err := generateFilePathAndName(req, file.GetClientOriginalName())
	if err != nil {
		return Error(ctx, http.StatusBadRequest, err.Error())
	}

	dokumen, err := facades.Storage().PutFileAs(path, file, nameLaporan)
	if err != nil {
		return ctx.Response().Status(422).Json(map[string]string{"error": "error writing file: " + err.Error()})
	}

	laporan := models.Laporan{
		NamaLaporan:  nameLaporan,
		JenisLaporan: req.JenisLaporan,
		TahunKe:      req.TahunKe,
		BulanKe:      req.BulanKe,
		MingguKe:     req.MingguKe,
		Dokumen:      dokumen,
	}

	var krywn, atasan models.Karyawan
	facades.Orm().Query().Where("emp_no = ?", req.CreatedBy).First(&krywn)
	if krywn.EmpNo == "" {
		return ErrorSystem(ctx, "Data Tidak Ditemukan")
	}
	facades.Orm().Query().Where("emp_no = ?", krywn.IDAtasan).First(&atasan)
	if atasan.EmpNo == "" {
		return ErrorSystem(ctx, "Data Tidak Ditemukan")
	}

	if err := facades.Orm().Query().Create(&laporan); err != nil {
		return ErrorSystem(ctx, "Data Gagal Ditambahkan")
	}

	approval := models.Approval{
		LaporanID:  laporan.IDLaporan,
		Status:     "WaitApproved",
		ApprovedBy: atasan.EmpNo,
	}

	facades.Orm().Query().Create(&approval)

	return ctx.Response().Json(http.StatusOK, map[string]string{"success": "success writing file"})
}

func (r *Laporan) Edit(ctx http.Context) http.Response {
	return nil
}

func (r *Laporan) Update(ctx http.Context) http.Response {
	return nil
}

func (r *Laporan) Destroy(ctx http.Context) http.Response {
	return nil
}

func (r *Laporan) ListLaporan(ctx http.Context) http.Response {
	// var laporan []models.Laporan

	// // Bind request data to req
	// if err := ctx.Request().Bind(&req); err != nil {
	// 	return ctx.Response().Json(http.StatusBadRequest, map[string]any{
	// 		"error": "Invalid request data",
	// 	})
	// }

	// // Correct the SQL syntax in Join and Where
	// facades.Orm().Query().Join("inner join public.approval apv on apv.laporan_id = id_laporan").
	// 	Where("(apv.status != ?) AND (apv.status != ?) AND ()", "WaitApproved", "Rejected").Join("").Group("id_laporan").Find(&laporan)

	// // Return the response
	// return ctx.Response().Json(http.StatusOK, map[string]any{
	// 	"laporan": laporan,
	// })
	return nil
}

func generateFilePathAndName(req lp.PostLaporan, originalFileName string) (string, string, error) {
	var path, nameLaporan string
	var jenis int

	switch {
	case req.BulanKe == 3 || req.BulanKe == 6 || req.BulanKe == 9:
		// Triwulan
		jenis = req.BulanKe / 3
		path = strconv.Itoa(req.TahunKe) + "/" + req.JenisLaporan + "/Bulan " + monthName(req.BulanKe)
		nameLaporan = "Laporan Triwulan ke-" + strconv.Itoa(jenis) + " " + strconv.Itoa(req.TahunKe)
	case req.MingguKe != 0 && req.BulanKe != 0 && req.TahunKe != 0:
		// Mingguan
		path = strconv.Itoa(req.TahunKe) + "/" + req.JenisLaporan + "/Bulan " + monthName(req.BulanKe)
		nameLaporan = "Laporan Minggu ke-" + strconv.Itoa(req.MingguKe) + " " + monthName(req.BulanKe) + " " + strconv.Itoa(req.TahunKe)
	case req.BulanKe != 0 && req.TahunKe != 0:
		// Bulanan
		path = strconv.Itoa(req.TahunKe) + "/" + req.JenisLaporan + "/Bulan " + monthName(req.BulanKe)
		nameLaporan = "Laporan Bulan " + monthName(req.BulanKe)
	default:
		return "", "", errors.New("Invalid request parameters")
	}

	return path, nameLaporan, nil
}

func monthName(month int) string {
	months := [12]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli",
		"Agustus", "September", "Oktober", "November", "Desember"}

	return months[month-1]
}
