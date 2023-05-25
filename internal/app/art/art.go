// Package art provides functionalities related to Art and Variant of the game,
// including serving them, manipulating them and keep them persisted.
//
// This package exposes its main interfaces by fx via FxModule function. However,
// it is still easy to instantiate an art.Service by simply creating it.
package art

import (
	"go.uber.org/fx"
)

// FxModule teaches a fx.App how to instantiate an art.Service, and also
// optionally register its HTTP endpoints on api.V1 if provided.
func FxModule() fx.Option {
	return fx.Module("art",
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
