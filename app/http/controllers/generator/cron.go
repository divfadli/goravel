package generator

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
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
func WeeklyGenerate() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	defer func() { _ = s.Shutdown() }()

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
		)),
		gocron.NewTask(
			// TODO: generate laporan mingguan
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
