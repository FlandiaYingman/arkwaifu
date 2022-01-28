package arkres

import (
	"context"
	"encoding/json"
	"github.com/facette/natsort"
	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	githubArknightsGameDataOwner = "Kengxxiao"
	githubArknightsGameDataRepo  = "ArknightsGameData"
)

func getGameDataJson(resVersion string, fullpath string) (string, error) {
	commitSha, err := getCommitShaByResVersion(resVersion)
	if err != nil {
		return "", err
	}

	client := github.NewClient(nil)
	fileContent, _, _, err := client.Repositories.GetContents(
		context.Background(),
		githubArknightsGameDataOwner,
		githubArknightsGameDataRepo,
		fullpath,
		&github.RepositoryContentGetOptions{
			Ref: commitSha,
		},
	)
	if err != nil {
		return "", err
	}

	if fileContent != nil {
		return fileContent.GetContent()
	} else {
		return "", errors.Errorf("%v is not a file or not exist", fullpath)
	}
}

func getCommitShaByResVersion(resVersion string) (string, error) {
	client := github.NewClient(nil)
	page := 1
	perPage := 100
	for true {
		commits, _, err := client.Repositories.ListCommits(
			context.Background(),
			githubArknightsGameDataOwner,
			githubArknightsGameDataRepo,
			&github.CommitsListOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: perPage,
				},
			},
		)
		if err != nil {
			return "", err
		}

		for _, c := range commits {
			message := c.GetCommit().GetMessage()
			if strings.Contains(message, "CN UPDATE") && strings.Contains(message, resVersion) {
				return c.GetCommit().GetSHA(), nil
			}
		}

		if len(commits) == 0 {
			break
		}
		page += 1
	}
	return "", errors.Errorf("commit by res version %v not found", resVersion)
}

func GetStoryReviewData(resVersion string) ([]StoryReviewData, error) {
	path := "zh_CN/gamedata/excel/story_review_table.json"
	gamedataJson, err := getGameDataJson(resVersion, path)
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
	StoryID    string  `json:"storyId"`
	StoryGroup string  `json:"storyGroup"`
	StorySort  int64   `json:"storySort"`
	StoryCode  *string `json:"storyCode"`
	StoryName  string  `json:"storyName"`
	StoryTxt   string  `json:"storyTxt"`
	AvgTag     AvgTag  `json:"avgTag"`
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

// GetStoryTxt gets the story txt
func GetStoryTxt(resVersion string, name string) (string, error) {
	infos, err := GetResInfos(resVersion)
	if err != nil {
		return "", err
	}
	err = GetTextRes(FilterResInfosRegexp(infos, regexp.MustCompile("gamedata/story")), "tmp")
	if err != nil {
		return "", err
	}

	fileBytes, err := ioutil.ReadFile(filepath.Join("tmp", "assets/torappu/dynamicassets", name))
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}
