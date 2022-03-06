package updateloop

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/gamedata"
	"os"
	"strings"
)

func GetAvgGameData(resVersion string) ([]avg.Group, error) {
	tempDir, err := os.MkdirTemp("", "arkwaifu-updateloop-avg_gamedata-*")
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
		linq.From(images).
			Select(func(i interface{}) interface{} { return strings.ToLower(i.(string)) }).
			ToSlice(&images)
		linq.From(backgrounds).
			Select(func(i interface{}) interface{} { return strings.ToLower(i.(string)) }).
			ToSlice(&backgrounds)
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
