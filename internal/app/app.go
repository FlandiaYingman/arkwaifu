package app

import (
	"arkwaifu/internal/app/config"
	"arkwaifu/internal/app/infra"
	repo2 "arkwaifu/internal/app/repo"
	"context"
	"fmt"
	"go.uber.org/fx"
)

func ProvideOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(config.ProvideConfig),
		fx.Provide(infra.ProvidePostgres),
		fx.Provide(
			repo2.NewAvgRepo,
			repo2.NewAvgGroupRepo,
		),
		//...
	}
	return opts
}

func Run(avgRepo *repo2.AvgRepo, avgGroupRepo *repo2.AvgGroupRepo) {
	id, err := avgGroupRepo.GetAvgGroupByID(context.Background(), "act15side")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *id)
	for _, avg := range id.Avgs {
		fmt.Printf("%+v\n", *avg)
	}
}
