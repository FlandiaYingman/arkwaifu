package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkdata"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkparser"
	"github.com/rs/zerolog/log"
	_ "image/jpeg" // register jpeg codec
	_ "image/png"  // register png codec
	"os"
	"time"
)

var (
	storyPatterns = []string{
		"gamedata/excel/**",
		"gamedata/story/**",
	}
)

func (s *Service) getRemoteStoryVersion(ctx context.Context, server ark.Server) (ark.Version, error) {
	resourceVersion, err := arkdata.GetLatestDataVersion(ctx, server)
	if err != nil {
		return "", err
	}
	return resourceVersion.ResourceVersion, err
}

func (s *Service) getLocalStoryVersion(_ context.Context, server ark.Server) (ark.Version, error) {
	return s.repo.selectStoryVersion(server)
}

func (s *Service) attemptUpdateStory(ctx context.Context, server ark.Server) {
	log := log.With().
		Str("server", server).
		Logger()

	log.Info().
		Msg("Attempting to update stories of the server...")

	localStoryVersion, err := s.getLocalStoryVersion(ctx, server)
	if err != nil {
		log.Err(err).
			Msg("Failed to get the local story version of the server. ")
		return
	}

	remoteStoryVersion, err := s.getRemoteStoryVersion(ctx, server)
	if err != nil {
		log.Err(err).
			Msg("Failed to get the remote story version of the server. ")
		return
	}

	log = log.With().
		Str("localStoryVersion", localStoryVersion).
		Str("remoteStoryVersion", remoteStoryVersion).
		Logger()

	if localStoryVersion != remoteStoryVersion {
		log.Info().
			Msg("Updating the stories of the server, since the story versions are not identical. ")
		begin := time.Now()
		err = s.updateStories(ctx, server, remoteStoryVersion)
		if err != nil {
			log.Err(err).
				Msg("Update loop failed to update the stories")
		}
		log.Info().
			Dur("elapsed", time.Since(begin)).
			Msg("Updated the stories of the server successfully. ")
	} else {
		log.Info().
			Msg("Skip updating the stories of the server, since the story versions are identical. ")
	}
}

func (s *Service) updateStories(ctx context.Context, server ark.Server, version ark.Version) error {
	root, err := os.MkdirTemp("", "arkwaifu-updateloop-story-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(root)

	err = arkdata.GetGameData(ctx, server, version, storyPatterns, root)
	if err != nil {
		return err
	}

	parser := arkparser.Parser{
		Root:   root,
		Prefix: "assets/torappu/dynamicassets",
	}

	storyTree, err := parser.Parse()
	if err != nil {
		return err
	}
	err = s.storyService.PopulateFrom(storyTree, server)
	if err != nil {
		return err
	}

	err = s.repo.upsertStoryVersion(&storyVersion{
		Server:  server,
		Version: version,
	})
	if err != nil {
		return err
	}

	return nil
}
