package generator

import (
	"fmt"
	"goravel/app/models"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/goravel/framework/facades"
)

func safeExecute(job func()) {
	defer func() {
		if r := recover(); r != nil {
			facades.Log().Error(fmt.Sprintf("Job panic recovered: %v", r))
		}
	}()
	job()
}

func QuarterlyGenerate() {
	s, _ := gocron.NewScheduler()

	// Create jobs and store their IDs
	jobs := make(map[string]uuid.UUID)

	// QuarterlyJob Every 3 Month
	QuarterlyJob, _ := s.NewJob(
		gocron.MonthlyJob(3, gocron.NewDaysOfTheMonth(-1), gocron.NewAtTimes(
			gocron.NewAtTime(23, 59, 0),
		)),
		gocron.NewTask(
			func() {
				safeExecute(func() {
					currentTime := time.Now()
					month := currentTime.Month()

					// Only generate for quarter-end months (3,6,9,12)
					if month == time.March || month == time.June ||
						month == time.September || month == time.December {
						generate := NewPdf("")
						generate.LaporanTriwulan()
					}
				})
			},
		),
		gocron.WithName("QuarterlyReportJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Quarterly Job %s with ID %s is about to run\n", jobName, jobID)
			}),
		),
	)
	jobs["QuarterlyReportJob"] = QuarterlyJob.ID()

	// Update QuarterlyUpdateJob Every Saturday
	QuarterlyUpdateJob, _ := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
		)),
		gocron.NewTask(
			// generate laporan quarterly
			func() {
				safeExecute(func() {
					generate := NewPdf("")
					var laporan []models.Laporan
					facades.Orm().Query().
						Join("INNER JOIN public.approval apv ON apv.laporan_id = laporan.id_laporan").
						Where("apv.status = ? AND laporan.jenis_laporan = ?", "Rejected", "Laporan Triwulan").
						Find(&laporan)

					for _, lp := range laporan {
						fmt.Println(lp.IDLaporan)
						generate.LaporanTriwulanUpdate(lp.IDLaporan,
							lp.BulanKe, lp.TahunKe)
					}
				})
			},
		),
		gocron.WithName("QuarterlyReportUpdateJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)
	jobs["QuarterlyReportUpdateJob"] = QuarterlyUpdateJob.ID()

	// Display all job IDs
	fmt.Println("=== All Scheduled Quarterly Jobs ===")
	for name, id := range jobs {
		fmt.Printf("Job Name: %s, ID: %s\n", name, id)
	}

	s.Start()
}

// https://pkg.go.dev/github.com/go-co-op/gocron/v2#MonthlyJob
func MonthlyGenerate() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// defer func() { _ = s.Shutdown() }()
	// Create jobs and store their IDs
	jobs := make(map[string]uuid.UUID)

	// add a job to the scheduler
	monthlyJob, _ := s.NewJob(
		gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(-1), gocron.NewAtTimes(
			gocron.NewAtTime(23, 59, 00),
		)),
		gocron.NewTask(
			// TODO: generate laporan bulanan
			func() {
				safeExecute(func() {
					generate := NewPdf("")
					generate.LaporanBulanan()
				})
			},
		),
		gocron.WithName("MonthlyReportJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)
	jobs["MonthlyReportJob"] = monthlyJob.ID()

	// Update MonthlyReportJob every Saturday
	monthlyUpdateJob, _ := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
		)),
		gocron.NewTask(
			// generate laporan Bulanan
			func() {
				safeExecute(func() {
					generate := NewPdf("")
					var laporan []models.Laporan
					facades.Orm().Query().
						Join("INNER JOIN public.approval apv ON apv.laporan_id = laporan.id_laporan").
						Where("apv.status = ? AND laporan.jenis_laporan = ?", "Rejected", "Laporan Bulanan").
						Find(&laporan)

					for _, lp := range laporan {
						fmt.Println(lp.IDLaporan)
						generate.LaporanBulananUpdate(lp.IDLaporan,
							lp.BulanKe, lp.TahunKe)
					}
				})
			},
		),
		gocron.WithName("MonthlyReportUpdateJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)
	jobs["MonthlyReportUpdateJob"] = monthlyUpdateJob.ID()

	// Display all job IDs
	fmt.Println("=== All Scheduled Monthly Jobs ===")
	for name, id := range jobs {
		fmt.Printf("Job Name: %s, ID: %s\n", name, id)
	}

	// start the scheduler
	s.Start()
}

// https://pkg.go.dev/github.com/go-co-op/gocron/v2#WeeklyJob
func WeeklyGenerate() {
	// create a scheduler
	s, _ := gocron.NewScheduler()

	// defer func() { _ = s.Shutdown() }()

	// Create jobs and store their IDs
	jobs := make(map[string]uuid.UUID)

	// add a job to the scheduler
	weeklyJob, _ := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			gocron.NewAtTime(00, 00, 00),
		)),
		gocron.NewTask(
			// generate laporan mingguan
			func() {
				safeExecute(func() {
					generate := NewPdf("")
					generate.LaporanMingguan()
				})
			},
		),
		gocron.WithName("WeeklyReportJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)
	jobs["WeeklyReportJob"] = weeklyJob.ID()

	// Update WeeklyReportJob
	weeklyUpdateJob, _ := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			// gocron.NewAtTime(00, 00, 00),
			gocron.NewAtTime(0, 0, 00),
		)),
		gocron.NewTask(
			// generate laporan mingguan
			func() {
				safeExecute(func() {
					generate := NewPdf("")

					var laporan []models.Laporan
					facades.Orm().Query().
						Join("INNER JOIN public.approval apv ON apv.laporan_id = laporan.id_laporan").
						Where("apv.status = ? AND laporan.jenis_laporan = ?", "Rejected", "Laporan Mingguan").
						Find(&laporan)
					for _, lp := range laporan {
						generate.LaporanMingguanUpdate(lp.IDLaporan, lp.MingguKe, lp.BulanKe, lp.TahunKe)
					}
				})
			},
		),
		gocron.WithName("WeeklyReportUpdateJob"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)
	jobs["WeeklyReportUpdateJob"] = weeklyUpdateJob.ID()

	// Display all job IDs
	fmt.Println("=== All Scheduled Weekly Jobs ===")
	for name, id := range jobs {
		fmt.Printf("Job Name: %s, ID: %s\n", name, id)
	}

	// start the scheduler
	s.Start()
}

// https://pkg.go.dev/github.com/go-co-op/gocron/v2#CronJob
func TestGenerateCronTab() {
	// create a scheduler
	s, _ := gocron.NewScheduler()
	// if err != nil {
	// 	// handle error
	// 	log.Fatal(err)
	// }

	// defer func() { _ = s.Shutdown() }()

	j, err := s.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			"* * * * *",
			false,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
				x := NewPdf("")
				x.LaporanMingguanUpdate(4, 1, 9, 2024)
			},
			"hello",
			1,
		),
		gocron.WithName("TestGenerateCronTab"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s is about to run\n", jobName, jobID)
			}),
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				fmt.Printf("Job %s with ID %s has run\n", jobName, jobID)
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				fmt.Printf("Job %s with ID %s has run with error %v\n", jobName, jobID, err)
			}),
		),
	)

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// each job has a unique id
	fmt.Println("Job TestGenerateCronTab ID:", j.ID())

	// start the scheduler
	s.Start()
}

func StartCronJob() {
	// go TestGenerateCronTab()
	go WeeklyGenerate()
	go MonthlyGenerate()
	go QuarterlyGenerate()
}
