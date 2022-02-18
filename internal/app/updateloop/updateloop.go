package updateloop

import (
	"arkwaifu/internal/app/entity"
	"arkwaifu/internal/app/service"
	"arkwaifu/internal/pkg/arkres/gamedata"
	"arkwaifu/internal/pkg/arkres/resource"
	"context"
	"io/ioutil"
	"regexp"
)

type UpdateLoopController struct {
	avgService *service.AvgService
}

func NewUpdateLoopController(avgService *service.AvgService) *UpdateLoopController {
	return &UpdateLoopController{avgService: avgService}
}

func (c *UpdateLoopController) UpdateResources() error {
	ctx := context.Background()

	resVersion, err := GetLatestResVersion()
	if err != nil {
		return err
	}
	currentResVersion, err := c.avgService.GetResVersion(ctx)
	if err != nil {
		return err
	}

	// Test whether the resource is up-to-date.
	if resVersion == currentResVersion {
		return nil
	}

	tempDir, err := ioutil.TempDir("", "tmp-*")
	if err != nil {
		return err
	}

	avgGameData, err := GetAvgGameData(resVersion, tempDir)
	if err != nil {
		return err
	}

	//err = GetAvgResources(resVersion, tempDir)
	//if err != nil {
	//	return err
	//}

	return c.avgService.UpsertAvgs(context.Background(), resVersion, avgGameData)
}

func GetLatestResVersion() (string, error) {
	return resource.GetResVersion()
}

func GetAvgGameData(resVersion string, tempDir string) ([]entity.AvgGroup, error) {
	err := gamedata.Get(resVersion, "", tempDir)
	if err != nil {
		return nil, err
	}
	raw, err := gamedata.GetStoryReviewData(tempDir)
	if err != nil {
		return nil, err
	}

	return convertRawToAvgGroup(raw, tempDir)
}

func GetAvgResources(resVersion string, dest string) error {
	infos, err := resource.GetResInfos(resVersion)
	if err != nil {
		return err
	}

	infos = resource.FilterResInfosRegexp(infos, regexp.MustCompile("^avg/"))
	return resource.GetRes(infos, dest)
}

func convertRawToAvgGroup(raw []gamedata.StoryReviewData, gamedataDir string) ([]entity.AvgGroup, error) {
	groups := make([]entity.AvgGroup, len(raw))
	for i, d := range raw {
		avgs, err := convertRawToAvg(d.InfoUnlockDatas, gamedataDir)
		if err != nil {
			return nil, err
		}
		groups[i] = entity.AvgGroup{
			ID:   d.ID,
			Name: d.Name,
			Avgs: avgs,
		}
	}
	return groups, nil
}

func convertRawToAvg(raw []gamedata.StoryData, gamedataDir string) ([]*entity.Avg, error) {
	avgs := make([]*entity.Avg, len(raw))
	for i, data := range raw {
		text, err := gamedata.GetStoryText(gamedataDir, data.StoryTxt)
		if err != nil {
			return nil, err
		}
		avgs[i] = &entity.Avg{
			StoryID:   data.StoryID,
			StoryCode: data.StoryCode,
			StoryName: data.StoryName,
			StoryTxt:  text,
			AvgTag:    string(data.AvgTag),
			GroupID:   data.StoryGroup,
		}
	}
	return avgs, nil
}
