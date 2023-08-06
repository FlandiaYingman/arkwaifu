package gallery

import "go.uber.org/fx"

// FxModule teaches a fx.App how to instantiate an art.Service, and also
// optionally register its HTTP endpoints on api.V1 if provided.
func FxModule() fx.Option {
	return fx.Module("gallery",
		fx.Provide(
			newRepository,
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
