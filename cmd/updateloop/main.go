package main

import (
	"arkwaifu/internal/app"
	"arkwaifu/internal/app/updateloop"
	"context"
	"go.uber.org/fx"
)

func main() {
	var options []fx.Option
	options = append(options, app.ProvideOptions()...)
	options = append(options, fx.Invoke(run))

	fxApp := fx.New(options...)
	fxApp.Run()

	err := fxApp.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

func run(controller *updateloop.Controller) {
	err := controller.UpdateResources()
	if err != nil {
		panic(err)
	}
}
