// Package main - cronjob entry point
package main

import (
	"github.com/neurochar/backend/internal/app"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/app/fxboot"
	"go.uber.org/fx"
)

func main() {
	cfg := config.LoadConfig("configs/base.yml", "configs/base.local.yml")

	appOptions := fxboot.CronjobAppGetOptionsMap(app.IDCronjob, cfg)

	app := fx.New(
		fxboot.OptionsMapToSlice(appOptions)...,
	)

	app.Run()
}
