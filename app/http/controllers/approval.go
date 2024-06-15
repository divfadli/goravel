package controllers

import (
	"goravel/app/http/requests/approval"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type Approval struct {
	//Dependent services
}

func NewApproval() *Approval {
	return &Approval{
		//Inject services
	}
}

func (r *Approval) Index(ctx http.Context) http.Response {
	return nil
}

func (r *Approval) ListApproval(ctx http.Context) http.Response {
	var req approval.GetApproval

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var laporan []models.Laporan
	facades.Orm().Query().Join("inner join public.approval apv ON apv.laporan_id = id_laporan").
		Where("apv.status = ? AND apv.approved_by = ?", "WaitApproved", req.Nik).Find(&laporan)

	return Success(ctx, http.Json{
		"data_laporan": laporan,
	})
}

func (r *Approval) StoreApproval(ctx http.Context) http.Response {
	var req approval.PostApproval

	if chekRequestErr := ctx.Request().Bind(&req); chekRequestErr != nil {
		return SanitizeGet(ctx, chekRequestErr)
	}

	var approval models.Approval
	if err := facades.Orm().Query().
		Where("laporan_id =? AND approved_by =? AND status=?", req.IdLaporan, req.Nik, "WaitApproved").
		First(&approval); err != nil || approval.IDApproval == 0 {
		return Error(ctx, http.StatusNotFound, "Data tidak ditemukan")
	}

	var pesan string

	if req.Status == "Approved" {
		approval.Status = "Approved"

		var atasan, newAtasan models.Karyawan
		facades.Orm().Query().Where("emp_no =?", approval.ApprovedBy).With("Jabatan").First(&atasan)

		if atasan.Jabatan.Name != "Kepala Bakamla" {
			facades.Orm().Query().Where("emp_no =?", atasan.IDAtasan).First(&newAtasan)
		}

		newApproval := models.Approval{
			LaporanID:  approval.LaporanID,
			Status:     "WaitApproved",
			ApprovedBy: newAtasan.EmpNo,
		}

		facades.Orm().Query().Save(&approval)
		facades.Orm().Query().Create(&newApproval)
		pesan = "Data Berhasil Disetujui"
	} else if req.Status == "Rejected" {
		approval.Status = "Rejected"
		approval.Keterangan = &req.Keterangan

		facades.Orm().Query().Save(&approval)
		pesan = "Data Berhasil Ditolak"
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{"success": pesan})
}
