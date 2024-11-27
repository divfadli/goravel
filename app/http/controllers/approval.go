package controllers

import (
	"fmt"
	"goravel/app/http/controllers/generator"
	"goravel/app/http/requests/approval"
	"goravel/app/models"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/lib/pq"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Approval struct {
	pdf *generator.Pdf
	//Dependent services
}

func NewApproval() *Approval {
	return &Approval{
		pdf: generator.NewPdf(""),
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
	if err := facades.Orm().Query().With("Laporan").With("Karyawan").With("Karyawan.Jabatan").
		With("Karyawan.User").With("Karyawan.User.Role").
		Where("laporan_id =? AND approved_by =? AND status=?", req.IdLaporan, req.Nik, "WaitApproved").
		First(&approval); err != nil || approval.IDApproval == 0 {
		return Error(ctx, http.StatusNotFound, "Data tidak ditemukan")
	}

	var pesan string

	if req.Status == "Approved" {
		approval.Status = "Approved"
		go func() {
			if approval.Karyawan.Jabatan.Name == "Deputi Informasi, Hukum dan Kerja Sama" {
				sourcePath := "storage/app/" + approval.Laporan.Dokumen
				filenameWithoutExtension := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))

				OutputDir := "storage/temp/extracted/" + approval.Laporan.JenisLaporan + "/" + filenameWithoutExtension
				_ = os.MkdirAll(OutputDir, 0755)
				outputPathTemp := fmt.Sprintf("storage/temp/%s", approval.Laporan.JenisLaporan)
				_ = os.MkdirAll(outputPathTemp, 0755)

				pageCount, _ := api.PageCountFile(sourcePath)
				conf := model.NewDefaultConfiguration()
				var extractedPaths []string
				for i := 1; i <= pageCount-1; i++ {
					selectedPage := []string{strconv.Itoa(i)}
					err := api.ExtractPagesFile(sourcePath, OutputDir, selectedPage, conf)
					if err == nil {

						filename := fmt.Sprintf("%s/%s_page_%d.pdf", OutputDir, filenameWithoutExtension, i)

						extractedPaths = append(extractedPaths, filename)
					}
				}

				baseURL := "http://" + ctx.Request().Host()
				now := time.Now()

				date := fmt.Sprintf("%d %s %d", now.Day(), generator.MonthNameIndonesia(now.Month()), now.Year())
				fmt.Printf("%d %s %d\n", now.Day(), generator.MonthNameIndonesia(now.Month()), now.Year())
				outputPath := fmt.Sprintf("storage/temp/%s/output-ttd-acc.pdf", approval.Laporan.JenisLaporan)

				if approval.Laporan.JenisLaporan == "Laporan Mingguan" {
					templatePath := "templates/ttd-mingguan.html"
					newTemplatePath := "ttd-mingguan.html"

					templateData := struct {
						BaseURL    string
						Tanggal    string
						Jabatan    string
						Nama       string
						IsApproved bool
						Ttd        *string
						Nik        string
					}{
						BaseURL:    baseURL,
						Tanggal:    date,
						Jabatan:    approval.Karyawan.Jabatan.Name,
						Nama:       approval.Karyawan.Name,
						IsApproved: true == true,
						Ttd:        approval.Karyawan.Ttd,
						Nik:        approval.Karyawan.EmpNo,
					}

					if err := r.pdf.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
						// Generate Image
						success, _ := r.pdf.GenerateSlide(outputPath, "Mingguan/")
						if success {
							extractedPaths = append(extractedPaths, outputPath)
						}
					} else {
						fmt.Printf("Error: %v\n", err)
						return
					}

					err := api.MergeCreateFile(extractedPaths, sourcePath, false, model.NewDefaultConfiguration())
					if err != nil {
						fmt.Println("Error merging PDF files:", err)
						return
					}
				} else {
					templatePath := "templates/ttd-word.html"
					newTemplatePath := "ttd-word.html"

					templateData := struct {
						BaseURL    string
						Tanggal    string
						Jabatan    string
						Nama       string
						IsApproved bool
						Ttd        *string
						Nik        string
					}{
						BaseURL:    baseURL,
						Tanggal:    date,
						Jabatan:    approval.Karyawan.Jabatan.Name,
						Nama:       approval.Karyawan.Name,
						IsApproved: true == true,
						Ttd:        approval.Karyawan.Ttd,
						Nik:        approval.Karyawan.EmpNo,
					}

					if err := r.pdf.ParseTemplate(templatePath, newTemplatePath, templateData); err == nil {
						success, _ := r.pdf.GenerateLaporan(outputPath, pageCount, approval.Laporan.JenisLaporan+"/")
						if success {
							extractedPaths = append(extractedPaths, outputPath)
						}
					} else {
						fmt.Println(err)
					}
					err := api.MergeCreateFile(extractedPaths, sourcePath, false, model.NewDefaultConfiguration())
					if err != nil {
						fmt.Println("Error merging PDF files:", err)
						return
					}
				}
				_ = os.RemoveAll(OutputDir)
				_ = os.RemoveAll(outputPathTemp)
			}
		}()
		var newAtasan models.Karyawan

		if approval.Karyawan.Jabatan.Name != "Kepala Bakamla" {
			facades.Orm().Query().Where("emp_no =?", approval.Karyawan.IDAtasan).First(&newAtasan)
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
		go func() {
			if approval.Laporan.JenisLaporan == "Laporan Mingguan" {
				r.ResetDataMingguan(approval.Laporan.IDLaporan, approval.Laporan.MingguKe, approval.Laporan.BulanKe, approval.Laporan.TahunKe)
			} else if approval.Laporan.JenisLaporan == "Laporan Bulanan" {
				r.ResetDataBulanan(approval.Laporan.IDLaporan, approval.Laporan.BulanKe, approval.Laporan.TahunKe)
			} else if approval.Laporan.JenisLaporan == "Laporan Triwulan" {
				r.ResetDataTriwulan(approval.Laporan.IDLaporan, approval.Laporan.BulanKe, approval.Laporan.TahunKe)
			}
		}()

		approval.Status = "Rejected"
		approval.Keterangan = &req.Keterangan

		facades.Orm().Query().Save(&approval)
		pesan = "Data Berhasil Ditolak"
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{"success": pesan})
}

func (r *Approval) ResetDataMingguan(id_laporan int64, weeksTo int, month int, years int) {
	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	dayperweek := 7

	jumlahHari := generator.DaysInMonth(now)

	var minggu []time.Time

	// Calculate the first day of the month
	var firstDayMonth, startOfMonth time.Time
	firstDayMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startOfMonth = firstDayMonth
	dayOfWeek := firstDayMonth.Weekday()

	daysToAdd := 0
	switch dayOfWeek {
	case time.Monday:
		daysToAdd = 6
	case time.Tuesday:
		daysToAdd = 5
	case time.Wednesday:
		daysToAdd = 4
	case time.Thursday:
		daysToAdd = 3
	case time.Friday:
		daysToAdd = 2
	case time.Saturday:
		daysToAdd = 8
	case time.Sunday:
		daysToAdd = 7
	}

	// minggu = append(minggu, firstDay)
	nextWeek := firstDayMonth.AddDate(0, 0, daysToAdd)
	minggu = append(minggu, nextWeek)

	jumlahHari -= daysToAdd + 1
	completeWeeks := jumlahHari / dayperweek
	remainingDays := jumlahHari % dayperweek

	// Iterate through each week and calculate the start date
	for i := 1; i <= completeWeeks; i++ {
		var startDate time.Time
		if i == completeWeeks {
			if remainingDays == 1 {
				startDate = nextWeek.AddDate(0, 0, (i*dayperweek)+remainingDays)
			} else {
				startDate = nextWeek.AddDate(0, 0, i*dayperweek)
			}
		} else {
			startDate = nextWeek.AddDate(0, 0, i*dayperweek)
		}
		minggu = append(minggu, startDate)
	}

	// Add the start date for the remaining days if any
	if remainingDays > 1 {
		startDate := nextWeek.AddDate(0, 0, (completeWeeks*dayperweek)+remainingDays)
		minggu = append(minggu, startDate)
	}

	var startOfWeek time.Time
	var endOfWeek time.Time
	for i := range minggu {
		if i == weeksTo {
			if i == 1 {
				startOfWeek = startOfMonth
				endOfWeek = minggu[i-1]
			} else {
				startOfWeek = minggu[i-2]
				endOfWeek = minggu[i-1]
			}
			break
		}
	}

	var weeklyDataKeamanans, weeklyDataKeselamatans generator.WeeklyData
	var tanggal = "tanggal > ? AND tanggal <=?"
	if startOfWeek.Day() == 1 {
		tanggal = "tanggal >= ? AND tanggal <=?"
	}
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)
	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where(tanggal, startOfWeek, endOfWeek).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	var result_keamanan []models.KejadianKeamananImage
	var result_keselamatan []models.KejadianKeselamatanImage
	var data_keamanan []models.KejadianKeamanan
	var data_keselamatan []models.KejadianKeselamatan
	// Now you can loop through the grouped data
	for _, week := range weeklyDataKeamanans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keamanan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

		for _, data := range data_keamanan {
			var data_keamanan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

			result_keamanan = append(result_keamanan, models.KejadianKeamananImage{
				KejadianKeamanan: data,
				FileImage:        data_keamanan_image,
			})
		}
	}
	sort.Slice(result_keamanan, func(i, j int) bool {
		if result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian == result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian {
			return result_keamanan[i].KejadianKeamanan.Tanggal.ToStdTime().Before(result_keamanan[j].KejadianKeamanan.Tanggal.ToStdTime())
		}
		return result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian < result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian
	})

	for _, week := range weeklyDataKeselamatans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keselamatan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)
		for _, data := range data_keselamatan {
			var data_keselamatan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

			result_keselamatan = append(result_keselamatan, models.KejadianKeselamatanImage{
				KejadianKeselamatan: data,
				FileImage:           data_keselamatan_image,
			})
		}
	}
	sort.Slice(result_keselamatan, func(i, j int) bool {
		if result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian == result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian {
			return result_keselamatan[i].KejadianKeselamatan.Tanggal.ToStdTime().Before(result_keselamatan[j].KejadianKeselamatan.Tanggal.ToStdTime())
		}
		return result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian < result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian
	})
	var id_keamanan_arr []int64
	for _, data := range result_keamanan {
		id_keamanan_arr = append(id_keamanan_arr, data.IdKejadianKeamanan)
	}
	var id_keselamatan_arr []int64
	for _, data := range result_keselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, data.IdKejadianKeselamatan)
	}
	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})
}
func (r *Approval) ResetDataBulanan(id_laporan int64, month int, years int) {
	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var weeklyDataKeamanans, weeklyDataKeselamatans generator.WeeklyData
	facades.Orm().Query().
		Table("public.kejadian_keamanan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keamanan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeamanans)

	facades.Orm().Query().
		Table("public.kejadian_keselamatan").
		Select("DATE_TRUNC('week', tanggal) AS week_start, ARRAY_AGG(id_kejadian_keselamatan) AS kejadian_ids").
		Where("tanggal >= ? AND tanggal <= ?", startOfMonth, endOfMonth).
		Group("DATE_TRUNC('week', tanggal)").
		Order("week_start asc").Get(&weeklyDataKeselamatans)

	var result_keamanan []models.KejadianKeamananImage
	var result_keselamatan []models.KejadianKeselamatanImage
	var data_keamanan []models.KejadianKeamanan
	var data_keselamatan []models.KejadianKeselamatan
	// Now you can loop through the grouped data
	for _, week := range weeklyDataKeamanans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keamanan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keamanan)

		for _, data := range data_keamanan {
			var data_keamanan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keamanan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keamanan_id=?", data.IdKejadianKeamanan).Find(&data_keamanan_image)

			result_keamanan = append(result_keamanan, models.KejadianKeamananImage{
				KejadianKeamanan: data,
				FileImage:        data_keamanan_image,
			})
		}
	}
	sort.Slice(result_keamanan, func(i, j int) bool {
		if result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian == result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian {
			return result_keamanan[i].KejadianKeamanan.Tanggal.ToStdTime().Before(result_keamanan[j].KejadianKeamanan.Tanggal.ToStdTime())
		}
		return result_keamanan[i].KejadianKeamanan.JenisKejadian.NamaKejadian < result_keamanan[j].KejadianKeamanan.JenisKejadian.NamaKejadian
	})

	for _, week := range weeklyDataKeselamatans {
		facades.Orm().Query().
			Join("inner join public.jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id ").
			With("JenisKejadian").Where("id_kejadian_keselamatan = ANY(?)", pq.Array(week.KejadianIDs)).
			Order("k.nama_kejadian asc, zona asc, tanggal asc").Find(&data_keselamatan)
		for _, data := range data_keselamatan {
			var data_keselamatan_image []models.FileImage
			facades.Orm().Query().Join("inner join public.image_keselamatan imk ON id_file_image = imk.file_image_id").
				Where("imk.kejadian_keselamatan_id=?", data.IdKejadianKeselamatan).Find(&data_keselamatan_image)

			result_keselamatan = append(result_keselamatan, models.KejadianKeselamatanImage{
				KejadianKeselamatan: data,
				FileImage:           data_keselamatan_image,
			})
		}
	}
	sort.Slice(result_keselamatan, func(i, j int) bool {
		if result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian == result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian {
			return result_keselamatan[i].KejadianKeselamatan.Tanggal.ToStdTime().Before(result_keselamatan[j].KejadianKeselamatan.Tanggal.ToStdTime())
		}
		return result_keselamatan[i].KejadianKeselamatan.JenisKejadian.NamaKejadian < result_keselamatan[j].KejadianKeselamatan.JenisKejadian.NamaKejadian
	})
	var id_keamanan_arr []int64
	for _, data := range result_keamanan {
		id_keamanan_arr = append(id_keamanan_arr, data.IdKejadianKeamanan)
	}
	var id_keselamatan_arr []int64
	for _, data := range result_keselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, data.IdKejadianKeselamatan)
	}

	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})
}
func (r *Approval) ResetDataTriwulan(id_laporan int64, month int, years int) {
	now := time.Date(years, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	year := strconv.Itoa(now.Year())

	quarters := map[time.Month]struct {
		quarter      string
		periodFormat string
		months       []string
	}{
		time.January:   {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.February:  {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.March:     {"I", "(JAN - MAR %s)", []string{"JAN", "FEB", "MAR"}},
		time.April:     {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.May:       {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.June:      {"II", "(APR - JUN %s)", []string{"APR", "MEI", "JUN"}},
		time.July:      {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.August:    {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.September: {"III", "(JUL - SEP %s)", []string{"JUL", "AGT", "SEP"}},
		time.October:   {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.November:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
		time.December:  {"IV", "(OCT - DEC %s)", []string{"OKT", "NOV", "DES"}},
	}

	quarterInfo := quarters[now.Month()]
	months := quarterInfo.months

	var dataKeamanan []models.KejadianKeamanan
	var dataKeselamatan []models.KejadianKeselamatan

	monthNumbers := make([]int, len(months))
	for i, month := range months {
		monthNumbers[i] = generator.MonthNameMap[month]
	}

	facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeamanan)
	facades.Orm().Query().Join("inner join jenis_kejadian k on k.id_jenis_kejadian = jenis_kejadian_id").
		With("JenisKejadian").
		Where("DATE_PART('month', tanggal) IN (?) AND DATE_PART('year', tanggal) = ?", monthNumbers, year).
		Order("k.nama_kejadian asc, tanggal asc").Find(&dataKeselamatan)

	var id_keamanan_arr []int64
	for _, kejadian := range dataKeamanan {
		id_keamanan_arr = append(id_keamanan_arr, kejadian.IdKejadianKeamanan)
	}
	var id_keselamatan_arr []int64
	for _, kejadian := range dataKeselamatan {
		id_keselamatan_arr = append(id_keselamatan_arr, kejadian.IdKejadianKeselamatan)
	}

	var kejadianKeamanan models.KejadianKeamanan
	facades.Orm().Query().Model(&kejadianKeamanan).Where("id_kejadian_keamanan IN (?)", id_keamanan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})

	var kejadianKeselamatan models.KejadianKeselamatan
	facades.Orm().Query().Model(&kejadianKeselamatan).Where("id_kejadian_keselamatan IN (?)", id_keselamatan_arr).Update(map[string]interface{}{
		"is_locked": false,
	})
}
