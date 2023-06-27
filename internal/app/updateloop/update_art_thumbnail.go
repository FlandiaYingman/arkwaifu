package updateloop

import (
	"bytes"
	"context"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"image"
	"runtime"
	"time"
)

func (s *Service) attemptUpdateArtThumbnails(ctx context.Context) {
	arts, err := s.artService.SelectArtsWhoseVariantAbsent("thumbnail")
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Failed to select arts.")
	}
	if len(arts) > 0 {
		err := s.updateArtsThumbnail(ctx, arts)
		if err != nil {
			log.Error().
				Err(err).
				Caller().
				Msg("Failed to update arts thumbnail.")
		}
	}
}

func (s *Service) updateArtsThumbnail(ctx context.Context, arts []*art.Art) error {
	log.Info().Msgf("Begin updating arts thumbnail, using %d goroutines", runtime.NumCPU())
	begin := time.Now()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())

	for _, art := range arts {
		art := art
		eg.Go(func() error {
			return s.updateArtThumbnail(ctx, art)
		})
	}

	err := eg.Wait()
	if err != nil {
		return err
	}

	log.Info().Msgf("End updating arts thumbnail, elapsed %v", time.Since(begin))
	return nil
}

func (s *Service) updateArtThumbnail(ctx context.Context, a *art.Art) error {
	content, err := s.artService.TakeContent(a.ID, art.VariationOrigin)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		return err
	}

	img = imaging.Fit(img, 240, 240, imaging.Lanczos)

	buf := bytes.Buffer{}
	err = webp.Encode(&buf, img, &webp.Options{
		Lossless: false,
		Quality:  50,
		Exact:    false,
	})
	if err != nil {
		return err
	}

	err = s.artService.UpsertVariants(art.NewVariant(a.ID, art.VariationThumbnail))
	if err != nil {
		return err
	}

	err = s.artService.StoreContent(a.ID, art.VariationThumbnail, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
