package app

import (
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/artext"
	"github.com/flandiayingman/arkwaifu/internal/app/gallery"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
	"go.uber.org/fx"
)

func FxModules() fx.Option {
	return fx.Options(
		artext.FxModule(),
		art.FxModule(),
		story.FxModule(),
		gallery.FxModule(),
	)
}
