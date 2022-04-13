package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flandiayingman/arkwaifu/internal/app"
	"github.com/flandiayingman/arkwaifu/internal/app/updateloop"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type updateloopTicker struct {
	ticker *time.Ticker
	done   chan struct{}
}

func newUpdateloopTicker() *updateloopTicker {
	return &updateloopTicker{
		ticker: time.NewTicker(5 * time.Minute),
		done:   make(chan struct{}),
	}
}
func (uc *updateloopTicker) close() {
	uc.ticker.Stop()
	uc.done <- struct{}{}
}

func main() {
	var options []fx.Option
	options = append(options, app.ProvideOptions()...)
	options = append(options,
		fx.Provide(newUpdateloopTicker),
		fx.Invoke(run),
	)

	fxApp := fx.New(options...)
	fxApp.Run()

	err := fxApp.Start(context.Background())
	if err != nil {
		panic(err)
	}
}

func run(lc fx.Lifecycle, ut *updateloopTicker, uc *updateloop.Service) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go updateResourcesLoop(ut, uc)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ut.close()
			return nil
		},
	})
}

func updateResourcesLoop(ut *updateloopTicker, uc *updateloop.Service) {
	// the ticker wouldn't emit a tick instantly, so update resources at first manually
	updateResources(uc)
	for {
		select {
		case <-ut.ticker.C:
			updateResources(uc)
		case <-ut.done:
			break
		}
	}
}

func updateResources(uc *updateloop.Service) {
	err := uc.AttemptUpdate(context.Background())
	if err != nil {
		log.WithField("error", fmt.Sprintf("%+v", err)).
			Error("error occurs during update resources")
	}
}
