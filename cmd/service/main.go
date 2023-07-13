package main

import (
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/artext"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
		fx.Provide(
			infra.ProvideConfig,
			infra.ProvideFiber,
			infra.ProvideGorm,
		),
		artext.FxModule(),
		art.FxModule(),
		story.FxModule(),
	)
	fxApp.Run()
}
