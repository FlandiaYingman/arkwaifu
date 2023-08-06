package main

import (
	"github.com/flandiayingman/arkwaifu/internal/app"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/app/updateloop"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
		fx.Provide(
			infra.ProvideConfig,
			infra.ProvideGorm,
		),
		app.FxModules(),
		updateloop.FxModule(),
	)
	fxApp.Run()
}
