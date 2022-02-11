package app

import (
	"arkwaifu/internal/config"
	"arkwaifu/internal/infra"
	"arkwaifu/internal/repo"
	"context"
	"fmt"
	"go.uber.org/fx"
)

func ProvideOptions() []fx.Option {
	opts := []fx.Option{
		fx.Provide(config.ProvideConfig),
		fx.Provide(infra.ProvidePostgres),
		fx.Provide(
			repo.NewAvgRepo,
			repo.NewAvgGroupRepo,
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
