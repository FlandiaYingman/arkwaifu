package updateloop

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alitto/pond"
	"github.com/chai2010/webp"
	"github.com/flandiayingman/arkwaifu/internal/app/art"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark/hgapi"
	"github.com/flandiayingman/arkwaifu/internal/pkg/cols"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"runtime"
	"time"
)

var (
	artRegexp = regexp.MustCompile("^avg/(imgs|bg|items|characters)")
	// artRegexp = regexp.MustCompile(`^avg/((imgs/avg_img_0_0\.ab)|(bg/avg_bkg_h1_bg_0_0\.ab)|(items/item_36_eu1\.ab)|(characters/(avg_npc_034\.ab|avg_4078_bdhkgt_1\.ab|avg_123_fang_1\.ab)))`)
)

type task func() error

func (s *Service) getRemoteArtVersion() (ark.Version, error) {
	return hgapi.GetResVersion()
}
func (s *Service) getLocalArtVersion() (ark.Version, error) {
	return s.repo.selectArtVersion()
}

func (s *Service) attemptUpdateArt(ctx context.Context) {
	localArtVersion, err := s.getLocalArtVersion()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Failed to get the remote art version.")
		return
	}
	remoteArtVersion, err := s.getRemoteArtVersion()
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("Failed to get the local art version.")
		return
	}
	if localArtVersion != remoteArtVersion {
		err := s.updateArts(ctx, localArtVersion, remoteArtVersion)
		if err != nil {
			log.Error().
				Err(err).
				Str("localArtVersion", localArtVersion).
				Str("remoteArtVersion", remoteArtVersion).
				Msg("Failed to update arts")
		}
	}
}

func (s *Service) updateArts(ctx context.Context, oldVersion, newVersion string) error {
	root, err := os.MkdirTemp("", "arkwaifu-updateloop-art-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(root) }()

	err = hgapi.GetFromHGAPI(context.Background(), oldVersion, newVersion, root, artRegexp)
	if err != nil {
		return err
	}

	pictureArtTasks, err := s.createPictureArtSubmitTasks(root)
	if err != nil {
		return err
	}
	characterArtTasks, err := s.createCharacterArtSubmitTasks(root)
	if err != nil {
		return err
	}

	workerNum := runtime.NumCPU()
	taskNum := len(pictureArtTasks) + len(characterArtTasks)
	log.Info().Msgf("Begin submitting arts, using %v workers to run %d tasks", workerNum, taskNum)

	begin := time.Now()
	pool := pond.New(workerNum, taskNum)
	defer pool.Stop()
	group, ctx := pool.GroupContext(ctx)
	for _, artTask := range pictureArtTasks {
		group.Submit(artTask)
	}
	for _, artTask := range characterArtTasks {
		group.Submit(artTask)
	}
	err = group.Wait()
	if err != nil {
		return err
	}

	err = s.repo.upsertArtVersion(&artVersion{
		Lock:    zeroPtr,
		Version: newVersion,
	})
	if err != nil {
		return err
	}

	log.Info().Msgf("End submitting arts, elapsed %v", time.Since(begin))
	return nil
}

func (s *Service) createPictureArtSubmitTasks(root string) (tasks []task, err error) {
	scanner := ark.Scanner{Root: root}

	pictureArts, err := scanner.ScanForPictureArts()
	if err != nil {
		return nil, err
	}

	for _, art := range pictureArts {
		tasks = append(tasks, s.createPictureArtSubmitTask(root, art))
	}

	return
}
func (s *Service) createPictureArtSubmitTask(root string, art *ark.PictureArt) task {
	return func() error {
		log.Info().Msgf("Submitting the picture art %v...", *art)
		err := s.submitPictureArt(root, art)
		if err != nil {
			return fmt.Errorf("failed to submit picture art %v: %w", art, err)
		}
		return nil
	}
}
func (s *Service) submitPictureArt(root string, pic *ark.PictureArt) error {
	processor := ark.Processor{Root: root}
	img, err := processor.ProcessPictureArt(pic)
	if err != nil {
		return err
	}

	err = s.artService.UpsertArts(art.NewArt(img.Art.ID, art.MustParseCategory(img.Art.Kind)))
	if err != nil {
		return err
	}

	err = s.artService.UpsertVariants(art.NewVariant(img.Art.ID, art.VariationOrigin))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = webp.Encode(&buf, img.Image, &webp.Options{Lossless: true})
	if err != nil {
		return err
	}
	err = s.artService.StoreContent(img.Art.ID, art.VariationOrigin, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createCharacterArtSubmitTasks(root string) (tasks []task, err error) {
	scanner := ark.Scanner{Root: root}

	characterArts, err := scanner.ScanForCharacterArts()
	if err != nil {
		return nil, err
	}

	for _, art := range characterArts {
		tasks = append(tasks, s.createCharacterArtSubmitTask(root, art))
	}

	return
}
func (s *Service) createCharacterArtSubmitTask(root string, art *ark.CharacterArt) task {
	return func() error {
		log.Info().Msgf("Submitting the character art %v...", *art)
		err := s.submitCharacterArt(root, art)
		if err != nil {
			return fmt.Errorf("failed to submit picture art %v: %w", art, err)
		}
		return nil
	}
}
func (s *Service) submitCharacterArt(root string, char *ark.CharacterArt) error {
	processor := ark.Processor{Root: root}
	imgs, err := processor.ProcessCharacterArt(char)
	if err != nil {
		return err
	}

	err = s.artService.UpsertArts(cols.Map(imgs, func(img ark.CharacterArtImage) *art.Art {
		return art.NewArt(img.ID(), art.MustParseCategory(img.Art.Kind))
	})...)
	if err != nil {
		return err
	}

	err = s.artService.UpsertVariants(cols.Map(imgs, func(img ark.CharacterArtImage) *art.Variant {
		return art.NewVariant(img.ID(), art.VariationOrigin)
	})...)
	if err != nil {
		return err
	}

	for _, img := range imgs {
		var buf bytes.Buffer
		err := webp.Encode(&buf, img.Image, &webp.Options{Lossless: true})
		if err != nil {
			return err
		}
		err = s.artService.StoreContent(img.ID(), art.VariationOrigin, buf.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}
