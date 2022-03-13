package updateloop

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/flandiayingman/arkwaifu/internal/app/avg"
	"github.com/flandiayingman/arkwaifu/internal/pkg/arkres/data"
	"os"
	"strings"
)

func GetAvgGameData(resVersion string) ([]avg.Group, error) {
	tempDir, err := os.MkdirTemp("", "arkwaifu-updateloop-avg_gamedata-*")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	err = data.Get(resVersion, "", tempDir)
	if err != nil {
		return nil, err
	}
	raw, err := data.GetStoryReviewData(tempDir)
	if err != nil {
		return nil, err
	}

	return groupsFromRaw(raw, tempDir)
}

func groupsFromRaw(raw []data.StoryReviewData, gamedataDir string) ([]avg.Group, error) {
	groups := make([]avg.Group, len(raw))
	for i, storyData := range raw {
		stories, err := storiesFromRaw(storyData, gamedataDir)
		if err != nil {
			return nil, err
		}
		groups[i] = avg.Group{
			ID:      storyData.ID,
			Name:    storyData.Name,
			ActType: string(storyData.ActType),
			Stories: stories,
		}
	}
	return groups, nil
}

func storiesFromRaw(reviewData data.StoryReviewData, gamedataDir string) ([]avg.Story, error) {
	raws := reviewData.InfoUnlockDatas
	stories := make([]avg.Story, len(raws))
	for i, raw := range raws {
		text, err := data.GetStoryText(gamedataDir, raw.StoryTxt)
		if err != nil {
			return nil, err
		}
		images, backgrounds := data.GetResourcesFromStoryText(text)
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
			GroupID:     reviewData.ID,
		}
	}
	return stories, nil
}
