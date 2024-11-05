package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// https://pkg.go.dev/github.com/go-co-op/gocron/v2#MonthlyJob
func MonthlyGenerate() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	defer func() { _ = s.Shutdown() }()

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.MonthlyJob(1, gocron.NewDaysOfTheMonth(-1), gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
		)),
		gocron.NewTask(
			// TODO: generate laporan bulanan
			func(a string, b int) {
				// do things
				fmt.Println(a, b)
			},
			"hello",
			1,
		),
	)

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()
}

// https://pkg.go.dev/github.com/go-co-op/gocron/v2#WeeklyJob
func CronJobGenerateLaporanMingguan() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// defer func() { _ = s.Shutdown() }()

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
		)),
		gocron.NewTask(
			// generate laporan mingguan
			func() {
				generate := NewPdf("")
				generate.LaporanMingguan()
			},
		),
		gocron.WithName("WeeklyGenerateLaporanMingguan"),
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
	fmt.Println("Job GenerateLaporanMingguan ID:", j.ID())

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
				x.LaporanMingguan()
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
	go CronJobGenerateLaporanMingguan()
}
