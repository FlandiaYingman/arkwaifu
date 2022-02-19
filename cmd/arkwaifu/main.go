package main

import (
	"arkwaifu/internal/app"
	"arkwaifu/internal/app/config"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"net"
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

func run(app *fiber.App, config *config.Config, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", config.Address)
			if err != nil {
				return err
			}
			go func() {
				if err := app.Listener(listener); err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
