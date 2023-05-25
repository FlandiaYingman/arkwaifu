package story

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Module("story",
		fx.Provide(
			newRepo,
			newService,
			newController,
		),
		fx.Invoke(
			fx.Annotate(
				registerController,
				fx.ParamTags(``, `optional:"true"`),
			),
		),
	)
}
