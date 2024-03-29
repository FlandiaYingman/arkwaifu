package main

import (
	"github.com/flandiayingman/arkwaifu/internal/app"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
		fx.Provide(
			infra.ProvideConfig,
			infra.ProvideFiber,
			infra.ProvideGorm,
		),
		app.FxModules(),
	)
	fxApp.Run()
}
