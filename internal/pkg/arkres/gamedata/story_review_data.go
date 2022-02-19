package gamedata

import (
	"encoding/json"
	"fmt"
	"github.com/facette/natsort"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	backgroundRegexp = regexp.MustCompile(`\[Background\(.*?image="(.*?)".*?\)]`)
	imageRegexp      = regexp.MustCompile(`\[Image\(.*?image="(.*?)".*?\)]`)
)

func GetStoryReviewData(gamedata string) ([]StoryReviewData, error) {
	dataPath := "excel/story_review_table.json"
	gamedataJson, err := GetText(gamedata, dataPath)
	if err != nil {
		return nil, err
	}

	var table StoryReviewTable
	err = json.Unmarshal([]byte(gamedataJson), &table)
	if err != nil {
		return nil, err
	}

	data := make([]StoryReviewData, 0, len(table))
	for _, v := range table {
		data = append(data, v)
	}

	sort.Sort(sortStoryReviewData(data))

	return data, nil
}

func GetStoryText(gamedataDir string, storyTextPath string) (string, error) {
	storyTextPath = fmt.Sprintf("%s.txt", storyTextPath)
	storyTextPath = filepath.Join("story", storyTextPath)
	return GetText(gamedataDir, storyTextPath)
}

func GetResourcesFromStoryText(storyText string) (images []string, backgrounds []string) {
	backgrounds = make([]string, 0)
	backgroundMatches := backgroundRegexp.FindAllStringSubmatch(storyText, -1)
	if backgroundMatches != nil {
		for _, match := range backgroundMatches {
			backgrounds = append(backgrounds, match[1])
		}
	}
	images = make([]string, 0)
	imageMatches := imageRegexp.FindAllStringSubmatch(storyText, -1)
	if imageMatches != nil {
		for _, match := range imageMatches {
			images = append(images, match[1])
		}
	}
	return
}

type StoryReviewTable map[string]StoryReviewData

type StoryReviewData struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	EntryType       EntryType   `json:"entryType"`
	ActType         ActType     `json:"actType"`
	StartTime       int64       `json:"startTime"`
	InfoUnlockDatas []StoryData `json:"infoUnlockDatas"`
}

type StoryData struct {
	StoryID    string `json:"storyId"`
	StoryGroup string `json:"storyGroup"`
	StorySort  int64  `json:"storySort"`
	StoryCode  string `json:"storyCode"`
	StoryName  string `json:"storyName"`
	StoryTxt   string `json:"storyTxt"`
	AvgTag     AvgTag `json:"avgTag"`
}

type ActType string

const (
	ActTypeNone   ActType = "NONE"
	ActivityStory ActType = "ACTIVITY_STORY"
	MainStory     ActType = "MAIN_STORY"
	MiniStory     ActType = "MINI_STORY"
)

type EntryType string

const (
	EntryTypeNone EntryType = "NONE"
	Activity      EntryType = "ACTIVITY"
	Mainline      EntryType = "MAINLINE"
	MiniActivity  EntryType = "MINI_ACTIVITY"
)

type AvgTag string

const (
	幕间  AvgTag = "幕间"
	行动前 AvgTag = "行动前"
	行动后 AvgTag = "行动后"
)

type sortStoryReviewData []StoryReviewData

func (s sortStoryReviewData) Len() int {
	return len(s)
}

func (s sortStoryReviewData) Less(i, j int) bool {
	if s[i].StartTime != s[j].StartTime {
		return s[i].StartTime < s[j].StartTime
	} else {
		return natsort.Compare(s[i].ID, s[j].ID)
	}
}

func (s sortStoryReviewData) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
