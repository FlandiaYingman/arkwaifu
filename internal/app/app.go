package app

import (
	"github.com/flandiayingman/arkwaifu/internal/app/asset"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/app/config"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/app/server"
	"github.com/flandiayingman/arkwaifu/internal/app/updateloop"
	"go.uber.org/fx"
)

func ProvideOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(config.ProvideConfig),
		fx.Provide(
			infra.ProvidePostgres,
			infra.ProvideFiber,
		),
		fx.Provide(server.ProvideV0),
		fx.Provide(
			avg.NewRepo,
			avg.NewService,
			avg.NewController,
		),
		fx.Invoke(
			avg.RegisterController,
		),
		fx.Provide(
			asset.NewRepo,
			asset.NewService,
			asset.NewController,
		),
		fx.Invoke(
			asset.RegisterController,
		),
		fx.Provide(
			updateloop.NewController,
		),
		// ...
	}
	return opts
}

func Run() {

}
