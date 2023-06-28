// Package updateloop provide functionalities to keep the local assets
// up-to-date.
//
// The operation "update" contains 3 sub-operations, and the subject of update
// can be either PictureArt module or Story module.
//
// Pull. It pulls the assets from the remote API, either the HyperGryph API or
// the ArknightsGameData GitHub repository API. It also unpacks or decompresses
// the assets from the downloaded files, therefore, whether the assets are
// pulled from either API, the result directory structure will turn out the
// same.
//
// Preprocess. It converts the pulled assets to the form which the submitter of
// the "submit" sub-operation can recognize. Moreover, it does some image
// preprocessing such as merge the alpha channel and merge the different faces
// onto the body for characters (合并差分).
//
// Submit. It parses the preprocessed assets and submits them to the art service
// and the story service.
//
// After the above sub-operations, all intermediate files shall be deleted to
// preserve the disk space.
package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/app/story"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"os"
	"time"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func FxModule() fx.Option {
	return fx.Module("updateloop",
		fx.Provide(
			newRepo,
			newService,
		),
		fx.Invoke(
			registerService,
		),
	)
}

// Service provides all functionalities of update.
type Service struct {
	repo *repo

	artService   *art.Service
	storyService *story.Service
}

func newService(
	artService *art.Service,
	storyService *story.Service,
	repo *repo,
) *Service {
	return &Service{
		artService:   artService,
		storyService: storyService,
		repo:         repo,
	}
}

func registerService(service *Service, lc fx.Lifecycle) {
	loopCtx, cancelLoop := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go service.Loop(loopCtx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancelLoop()
			return nil
		},
	})
}

func (s *Service) Loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			s.AttemptUpdate(ctx)
		}
		time.Sleep(5 * time.Minute)
	}
}

// AttemptUpdate attempts to update the resources.
func (s *Service) AttemptUpdate(ctx context.Context) {
	log := log.With().
		Logger()
	log.Info().Msg("Update loop is attempting to update the assets... ")

	s.attemptUpdateArt(ctx)
	s.attemptUpdateArtThumbnails(ctx)

	for _, server := range ark.Servers {
		// Skip TW server, since we haven't implemented it yet.
		if server == ark.TwServer {
			continue
		}
		s.attemptUpdateStory(ctx, server)
	}

	log.Info().Msg("Update loop has completed this attempt to update the assets..")
}
