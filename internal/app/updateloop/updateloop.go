package updateloop

import (
	"arkwaifu/internal/app/avg"
	"arkwaifu/internal/pkg/arkres/gamedata"
	"arkwaifu/internal/pkg/arkres/resource"
	"context"
	"io/ioutil"
	"regexp"
)

type Controller struct {
	avgService *avg.Service
}

func NewController(avgService *avg.Service) *Controller {
	return &Controller{avgService}
}

func (c *Controller) UpdateResources() error {
	ctx := context.Background()

	resVersion, err := GetLatestResVersion()
	if err != nil {
		return err
	}
	currentResVersion, err := c.avgService.GetVersion(ctx)
	if err != nil {
		return err
	}
	// Test whether the resource is up-to-date.
	if resVersion == currentResVersion {
		return nil
	}

	tempDir, err := ioutil.TempDir("", "arkwaifu-updateloop-*")
	if err != nil {
		return err
	}

	avgGameData, err := GetAvgGameData(resVersion, tempDir)
	if err != nil {
		return err
	}

	// err = GetAvgResources(resVersion, tempDir)
	// if err != nil {
	//	return err
	// }

	return c.avgService.SetAvgs(resVersion, avgGameData)
}

func GetLatestResVersion() (string, error) {
	return resource.GetResVersion()
}

func GetAvgGameData(resVersion string, tempDir string) ([]avg.Group, error) {
	err := gamedata.Get(resVersion, "", tempDir)
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
	return resource.GetRes(infos, dest)
}

func groupsFromRaw(raw []gamedata.StoryReviewData, gamedataDir string) ([]avg.Group, error) {
	groups := make([]avg.Group, len(raw))
	for i, data := range raw {
		stories, err := storiesFromRaw(data.InfoUnlockDatas, gamedataDir)
		if err != nil {
			return nil, err
		}
		groups[i] = avg.Group{
			ID:        data.ID,
			Name:      data.Name,
			StoryList: stories,
		}
	}
	return groups, nil
}

func storiesFromRaw(raw []gamedata.StoryData, gamedataDir string) ([]avg.Story, error) {
	stories := make([]avg.Story, len(raw))
	for i, data := range raw {
		text, err := gamedata.GetStoryText(gamedataDir, data.StoryTxt)
		if err != nil {
			return nil, err
		}
		images, backgrounds := gamedata.GetResourcesFromStoryText(text)
		stories[i] = avg.Story{
			ID:                data.StoryID,
			Code:              data.StoryCode,
			Name:              data.StoryName,
			Tag:               string(data.AvgTag),
			ImageResList:      images,
			BackgroundResList: backgrounds,
		}
	}
	return stories, nil
}
