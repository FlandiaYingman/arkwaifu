package app

import (
	"arkwaifu/internal/app/config"
	"arkwaifu/internal/app/infra"
	"arkwaifu/internal/app/repo"
	"arkwaifu/internal/app/service"
	"arkwaifu/internal/app/updateloop"
	"context"
	"fmt"
	"go.uber.org/fx"
)

func ProvideOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(config.ProvideConfig),
		fx.Provide(infra.ProvidePostgres),
		fx.Provide(
			repo.NewResVersionRepo,
			repo.NewAvgRepo,
			repo.NewAvgGroupRepo,
		),
		fx.Provide(
			service.NewAvgService,
		),
		fx.Provide(
			updateloop.NewUpdateLoopController,
		),
		//...
	}
	return opts
}

func Run(avgRepo *repo.AvgRepo, avgGroupRepo *repo.AvgGroupRepo) {
	id, err := avgGroupRepo.GetAvgGroupByID(context.Background(), "act15side")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *id)
	for _, avg := range id.Avgs {
		fmt.Printf("%+v\n", *avg)
	}
}
