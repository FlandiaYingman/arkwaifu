package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/infra"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"go.uber.org/fx"
	"testing"
)

func TestService_AttemptUpdate(t *testing.T) {
	fxApp := fx.New(
		fx.Options(
			fx.Provide(
				infra.ProvideConfig,
			),
			fx.Provide(
				infra.ProvideGorm,
				infra.ProvideGormNumericCollate,
				infra.ProvideFiber,
			),
			story.FxModule(),
			art.FxModule(),
			FxModule(),
		),
		fx.Invoke(func(service *Service, shutdowner fx.Shutdowner) {
			service.AttemptUpdate(context.Background(), ark.CnServer)
			err := shutdowner.Shutdown()
			if err != nil {
				panic(err)
			}
		}),
	)
	fxApp.Run()
}
