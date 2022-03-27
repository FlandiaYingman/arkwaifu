package arkavg

import (
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"path/filepath"
	"regexp"
)

type Avg struct {
	Groups []Group
}

type Group struct {
	ID      string  `bson:"id"`
	Name    string  `bson:"name"`
	Type    Type    `bson:"entryType"`
	Stories []Story `bson:"infoUnlockDatas"`
}

type Type string

const (
	TypeNone         Type = "NONE"
	TypeMain         Type = "MAINLINE"
	TypeActivity     Type = "ACTIVITY"
	TypeMiniActivity Type = "MINI_ACTIVITY"
)

type Story struct {
	ID      string `bson:"storyId"`
	GroupID string `bson:"storyGroup"`
	Code    string `bson:"storyCode"`
	Name    string `bson:"storyName"`
	Txt     string `bson:"storyTxt"`
	Tag     Tag    `bson:"avgTag"`
}

type Tag string

const (
	TagInterlude Tag = "幕间"
	TagBefore    Tag = "行动前"
	TagAfter     Tag = "行动后"
)

func GetAvg(resDir string, prefix string) (Avg, error) {
	bsonPath := "gamedata/excel/story_review_table.bson"
	bsonPath = filepath.Join(resDir, prefix, bsonPath)

	bsonContent, err := os.ReadFile(bsonPath)
	if err != nil {
		return Avg{}, errors.WithStack(err)
	}

	var ordered bson.D
	err = bson.Unmarshal(bsonContent, &ordered)
	if err != nil {
		return Avg{}, errors.WithStack(err)
	}

	b := bson.Raw(bsonContent)
	values, err := b.Values()
	groups := make([]Group, len(values))
	for i, e := range values {
		var group Group
		err := bson.Unmarshal(e.Value, &group)
		if err != nil {
			return Avg{}, errors.WithStack(err)
		}
		groups[i] = group
	}

	return Avg{Groups: groups}, nil
}

var (
	imageRegexp      = regexp.MustCompile(`(?i)\[Image\(.*?image="(.*?)".*?\)]`)
	backgroundRegexp = regexp.MustCompile(`(?i)\[Background\(.*?image="(.*?)".*?\)]`)
)

func GetStoryAssets(resDir string, prefix string, story Story) ([]Asset, error) {
	txt, err := GetStoryTxt(resDir, prefix, story)
	if err != nil {
		return nil, err
	}

	var assets []Asset
	assets = append(assets, findAssetsFromTxt(txt, KindImage, imageRegexp)...)
	assets = append(assets, findAssetsFromTxt(txt, KindBackground, backgroundRegexp)...)

	return assets, nil
}

func findAssetsFromTxt(txt string, kind Kind, regexp *regexp.Regexp) []Asset {
	var assets []Asset
	matches := regexp.FindAllStringSubmatch(txt, -1)
	for _, match := range matches {
		assets = append(assets, Asset{
			Name: match[1],
			Kind: kind,
		})
	}
	return assets
}

func GetStoryTxt(resDir string, prefix string, story Story) (string, error) {
	txtPath := fmt.Sprintf("gamedata/story/%v.txt", story.Txt)
	txtPath = filepath.Join(resDir, prefix, txtPath)

	txtContent, err := os.ReadFile(txtPath)
	if err != nil {
		return "", err
	}
	return string(txtContent), nil
}
