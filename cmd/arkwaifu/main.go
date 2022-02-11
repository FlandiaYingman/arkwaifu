package main

import (
	"arkwaifu/internal/app"
	"context"
	"go.uber.org/fx"
)

func main() {
	var options []fx.Option
	options = append(options, app.ProvideOptions()...)
	options = append(options, fx.Invoke(app.Run))

	fxApp := fx.New(options...)
	fxApp.Run()

	err := fxApp.Start(context.Background())
	if err != nil {
		return
	}
}
