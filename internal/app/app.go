package app

import (
	"arkwaifu/internal/app/avg"
	"arkwaifu/internal/app/config"
	"arkwaifu/internal/app/infra"
	"arkwaifu/internal/app/server"
	"arkwaifu/internal/app/updateloop"
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
			avg.NewVersionRepo,
			avg.NewStoryRepo,
			avg.NewGroupRepo,
			avg.NewService,
			avg.NewController,
		),
		fx.Invoke(
			avg.RegisterController,
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
