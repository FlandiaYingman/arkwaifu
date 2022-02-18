package gamedata

import (
	"encoding/json"
	"github.com/facette/natsort"
	"sort"
)

func GetStoryReviewData(gamedata string) ([]StoryReviewData, error) {
	dataPath := "zh_CN/gamedata/excel/story_review_table.json"
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
