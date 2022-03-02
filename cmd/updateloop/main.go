package main

import (
	"arkwaifu/internal/app"
	"arkwaifu/internal/app/updateloop"
	"context"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"time"
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

func run(lc fx.Lifecycle, ut *updateloopTicker, uc *updateloop.Controller) {
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

func updateResourcesLoop(ut *updateloopTicker, uc *updateloop.Controller) {
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

func updateResources(uc *updateloop.Controller) {
	err := uc.UpdateResources()
	if err != nil {
		log.WithError(err).Error("error occurs during update resources")
	}
}
