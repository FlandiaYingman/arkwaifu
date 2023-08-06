package app

import (
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/gallery"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
	"github.com/flandiayingman/arkwaifu/internal/app/updateloop"
	"go.uber.org/fx"
)

func FxOptions() fx.Option {
	return fx.Options(
		art.FxModule(),
		story.FxModule(),
		gallery.FxModule(),
		updateloop.FxModule(),
	)
}
