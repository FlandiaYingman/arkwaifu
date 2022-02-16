package updateloop

import (
	"arkwaifu/internal/app/entity"
	"arkwaifu/internal/app/service"
	"arkwaifu/internal/pkg/arkres"
	"context"
	"regexp"
)

type UpdateLoopController struct {
	avgService service.AvgService
}

func (c *UpdateLoopController) UpdateResources() error {
	resVersion, err := GetLatestResVersion()
	if err != nil {
		return err
	}

	avgData, err := GetAvgData(resVersion)
	if err != nil {
		return err
	}

	err = GetAvgResources(resVersion, "TODO/")
	if err != nil {
		return err
	}

	return c.avgService.UpsertAvgs(context.Background(), resVersion, avgData)
}

func GetLatestResVersion() (string, error) {
	return arkres.GetResVersion()
}

func GetAvgData(resVersion string) ([]entity.AvgGroup, error) {
	raw, err := arkres.GetStoryReviewData(resVersion)
	if err != nil {
		return nil, err
	}

	return convertRawToAvgGroup(raw), nil
}

func GetAvgResources(resVersion string, dest string) error {
	infos, err := arkres.GetResInfos(resVersion)
	if err != nil {
		return err
	}

	infos = arkres.FilterResInfosRegexp(infos, regexp.MustCompile("^avg/"))
	return arkres.GetRes(infos, dest)
}

func convertRawToAvgGroup(raw []arkres.StoryReviewData) []entity.AvgGroup {
	groups := make([]entity.AvgGroup, len(raw))
	for i, d := range raw {
		groups[i] = entity.AvgGroup{
			ID:   d.ID,
			Name: d.Name,
			Avgs: convertRawToAvg(d.InfoUnlockDatas),
		}
	}
	return groups
}

func convertRawToAvg(raw []arkres.StoryData) []*entity.Avg {
	avgs := make([]*entity.Avg, len(raw))
	for i, data := range raw {
		avgs[i] = &entity.Avg{
			StoryID:   data.StoryID,
			StoryCode: data.StoryCode,
			StoryName: data.StoryName,
			StoryTxt:  data.StoryTxt,
			AvgTag:    string(data.AvgTag),
			GroupID:   data.StoryGroup,
		}
	}
	return avgs
}
