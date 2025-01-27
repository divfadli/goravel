package main

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers/generator"
	"goravel/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	go func() {
		generator.StartCronJob()
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route run error: %v", err)
		}
	}()

	select {}
}
