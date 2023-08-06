package updateloop

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/app/gallery"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkdata"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkjson"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

var (
	galleryPatterns = []string{
		"gamedata/excel/story_review_meta_table.json",
		"gamedata/excel/replicate_table.json",
		"gamedata/excel/retro_table.json",
		"gamedata/excel/roguelike_topic_table.json",
	}
)

func (s *Service) getRemoteGalleryVersion(ctx context.Context, server ark.Server) (ark.Version, error) {
	v, err := arkdata.GetLatestDataVersion(ctx, server)
	return v.ResourceVersion, err
}

func (s *Service) getLocalGalleryVersion(ctx context.Context, server ark.Server) (ark.Version, error) {
	return s.repo.selectGalleryVersion(ctx, server)
}

func (s *Service) attemptUpdateGalleries(ctx context.Context, server ark.Server) {
	log := log.With().
		Str("server", server).
		Logger()
	log.Info().Msg("Attempting to update galleries of the server. ")
	localVersion, err := s.getLocalGalleryVersion(ctx, server)
	if err != nil {
		log.Err(err).
			Msg("Failed to get the local gallery version of the server. ")
		return
	}
	remoteVersion, err := s.getRemoteGalleryVersion(ctx, server)
	if err != nil {
		log.Err(err).
			Msg("Failed to get the remote gallery version of the server. ")
		return
	}
	log = log.With().
		Str("localVersion", localVersion).
		Str("remoteVersion", remoteVersion).
		Logger()
	if localVersion != remoteVersion {
		log.Info().Msg("Updating the stories of the server, since the gallery versions are not identical. ")
		begin := time.Now()
		err = s.updateGalleries(ctx, server, remoteVersion)
		if err != nil {
			log.Err(err).Msg("Update loop failed to update the stories")
		}
		log.Info().
			Str("elapsed", time.Since(begin).String()).
			Msg("Updated the stories of the server successfully. ")
	} else {
		log.Info().Msg("Skip updating the stories of the server, since the gallery versions are identical. ")
	}
}

func (s *Service) updateGalleries(ctx context.Context, server ark.Server, version ark.Version) error {
	root, err := os.MkdirTemp("", "arkwaifu-updateloop-gallery-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(root)

	err = arkdata.GetGameData(ctx, server, version, galleryPatterns, root)
	if err != nil {
		return err
	}

	galleries, err := ParseToGalleries(server, root)
	if err != nil {
		return err
	}

	err = s.galleryService.Put(galleries)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.repo.upsertGalleryVersion(ctx, &galleryVersion{
		Server:  server,
		Version: version,
	})
	if err != nil {
		return err
	}

	return nil
}

func ParseToGalleries(server ark.Server, root string) ([]gallery.Gallery, error) {
	jsonStoryReviewMetaTable, err := arkjson.Get(root, arkjson.StoryReviewMetaTablePath)
	if err != nil {
		return nil, err
	}
	jsonRetroTable, err := arkjson.Get(root, arkjson.RetroTable)
	if err != nil {
		return nil, err
	}
	jsonReplicateTable, err := arkjson.Get(root, arkjson.ReplicateTable)
	if err != nil {
		return nil, err
	}
	jsonRoguelikeTopicTable, err := arkjson.Get(root, arkjson.RoguelikeTopicTable)
	if err != nil {
		return nil, err
	}

	artMap := make(map[string]gallery.Art)
	for _, c := range jsonStoryReviewMetaTable.S("actArchiveResData", "pics").Children() {
		artMap[strings.ToLower(c.S("id").Data().(string))] = gallery.Art{
			Server:      server,
			GalleryID:   "", // Auto Generated
			SortID:      0,
			ID:          strings.ToLower(c.S("assetPath").Data().(string)),
			Name:        c.S("desc").Data().(string),
			Description: c.S("picDescription").Data().(string),
		}
	}

	galleryMap := make(map[string]gallery.Gallery)
	for _, c := range jsonRetroTable.S("retroActList").Children() {
		for _, actID := range c.S("linkedActId").Children() {
			galleryMap[strings.ToLower(actID.Data().(string))] = gallery.Gallery{
				Server:      server,
				ID:          strings.ToLower(actID.Data().(string)),
				Name:        c.S("name").Data().(string),
				Description: c.S("detail").Data().(string),
				Arts:        nil,
			}
		}
	}
	for _, c := range jsonRoguelikeTopicTable.S("topics").Children() {
		galleryMap[strings.ToLower(c.S("id").Data().(string))] = gallery.Gallery{
			Server:      server,
			ID:          strings.ToLower(c.S("id").Data().(string)),
			Name:        c.S("name").Data().(string),
			Description: c.S("lineText").Data().(string),
			Arts:        nil,
		}
	}

	galleries := make([]gallery.Gallery, 0)
	for id, c := range jsonStoryReviewMetaTable.Search("actArchiveData", "components").ChildrenMap() {
		if jsonReplicateTable.Exists(id) {
			continue
		}
		gallery := galleryMap[strings.ToLower(id)]
		for _, pic := range c.S("pic", "pics").Children() {
			art := artMap[strings.ToLower(pic.S("picId").Data().(string))]
			art.SortID = int(pic.S("picSortId").Data().(float64))
			gallery.Arts = append(gallery.Arts, art)
		}
		galleries = append(galleries, gallery)
	}

	return galleries, nil
}
