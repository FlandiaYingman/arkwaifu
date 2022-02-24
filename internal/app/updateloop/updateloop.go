package updateloop

import (
	"arkwaifu/internal/app/avg"
	"arkwaifu/internal/app/config"
	"arkwaifu/internal/pkg/arkres/gamedata"
	"arkwaifu/internal/pkg/arkres/resource"
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
)

type Controller struct {
	resLocation string
	avgService  *avg.Service
}

func NewController(avgService *avg.Service, conf *config.Config) *Controller {
	return &Controller{conf.ResourceLocation, avgService}
}

func (c *Controller) UpdateResources() error {
	ctx := context.Background()
	log.WithFields(log.Fields{}).Info("Attempt to update resources.")

	latestResVersion, err := GetLatestResVersion()
	if err != nil {
		return err
	}
	currentResVersion, err := c.avgService.GetVersion(ctx)
	if err != nil {
		return err
	}

	// Test whether the resource is up-to-date.
	logF := log.WithFields(log.Fields{
		"latestResVersion":  latestResVersion,
		"currentResVersion": currentResVersion,
	})
	if latestResVersion == currentResVersion {
		logF.Info("Resources are up-to-date.")
		return nil
	} else {
		logF.Info("Resource are out-of-date.")
	}

	logF.Info("Get AVG gamedata.")
	avgGameData, err := GetAvgGameData(latestResVersion)
	if err != nil {
		return err
	}

	resLocation := filepath.Join(c.resLocation, latestResVersion)
	logF.Info("Get AVG resources")
	err = GetAvgResources(latestResVersion, resLocation)
	if err != nil {
		return err
	}

	logF.Info("Set AVGs")
	return c.avgService.SetAvgs(latestResVersion, avgGameData)
}

func GetLatestResVersion() (string, error) {
	return resource.GetResVersion()
}

func GetAvgGameData(resVersion string) ([]avg.Group, error) {
	tempDir, err := os.MkdirTemp("", "arkwaifu-updateloop-*")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	err = gamedata.Get(resVersion, "", tempDir)
	if err != nil {
		return nil, err
	}
	raw, err := gamedata.GetStoryReviewData(tempDir)
	if err != nil {
		return nil, err
	}

	return groupsFromRaw(raw, tempDir)
}

func GetAvgResources(resVersion string, dest string) error {
	infos, err := resource.GetResInfos(resVersion)
	if err != nil {
		return err
	}

	infos = resource.FilterResInfosRegexp(infos, regexp.MustCompile("^avg/"))
	err = resource.GetRes(infos, dest)
	if err != nil {
		return err
	}

	err = os.Rename(filepath.Join(dest, "assets/torappu/dynamicassets/avg/images"), filepath.Join(dest, "images"))
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(dest, "assets/torappu/dynamicassets/avg/backgrounds"), filepath.Join(dest, "backgrounds"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(dest, "assets"))
	if err != nil {
		return err
	}
	return nil
}

func groupsFromRaw(raw []gamedata.StoryReviewData, gamedataDir string) ([]avg.Group, error) {
	groups := make([]avg.Group, len(raw))
	for i, data := range raw {
		stories, err := storiesFromRaw(data, gamedataDir)
		if err != nil {
			return nil, err
		}
		groups[i] = avg.Group{
			ID:      data.ID,
			Name:    data.Name,
			ActType: string(data.ActType),
			Stories: stories,
		}
	}
	return groups, nil
}

func storiesFromRaw(data gamedata.StoryReviewData, gamedataDir string) ([]avg.Story, error) {
	raws := data.InfoUnlockDatas
	stories := make([]avg.Story, len(raws))
	for i, raw := range raws {
		text, err := gamedata.GetStoryText(gamedataDir, raw.StoryTxt)
		if err != nil {
			return nil, err
		}
		images, backgrounds := gamedata.GetResourcesFromStoryText(text)
		stories[i] = avg.Story{
			ID:          raw.StoryID,
			Code:        raw.StoryCode,
			Name:        raw.StoryName,
			Tag:         string(raw.AvgTag),
			Images:      images,
			Backgrounds: backgrounds,
			GroupID:     data.ID,
		}
	}
	return stories, nil
}
