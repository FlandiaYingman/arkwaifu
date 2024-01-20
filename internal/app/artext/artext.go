package artext

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	registerController := func(controller *Controller, param RegisterParams) { controller.register(param) }

	return fx.Module("art-extension",
		fx.Provide(
			newService,
			newController,
		),
		fx.Invoke(registerController),
	)
}
