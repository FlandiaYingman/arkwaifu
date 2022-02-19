package app

import (
	"arkwaifu/internal/app/avg"
	"arkwaifu/internal/app/config"
	"arkwaifu/internal/app/infra"
	"arkwaifu/internal/app/updateloop"
	"go.uber.org/fx"
)

func ProvideOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(config.ProvideConfig),
		fx.Provide(infra.ProvidePostgres),
		fx.Provide(
			avg.NewVersionRepo,
			avg.NewStoryRepo,
			avg.NewGroupRepo,
			avg.NewService,
			avg.NewController,
		),
		fx.Provide(
			updateloop.NewController,
		),
		// ...
	}
	return opts
}

func Run(avgRepo *avg.StoryRepo, avgGroupRepo *avg.GroupRepo) {

}
