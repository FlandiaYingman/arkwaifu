package arkavg

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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

func GetStoryAssets(resDir string, prefix string, story Story) ([]Asset, error) {
	txt, err := GetStoryTxt(resDir, prefix, story)
	if err != nil {
		return nil, err
	}

	var assets []Asset
	assets = append(assets, findPicAssetsFromTxt(txt, KindImage, imageRegexp)...)
	assets = append(assets, findPicAssetsFromTxt(txt, KindBackground, backgroundRegexp)...)

	charAssets, err := findCharAssetsFromTxt(txt)
	if err != nil {
		return nil, err
	}
	assets = append(assets, charAssets...)

	return assets, nil
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

func findPicAssetsFromTxt(txt string, kind Kind, regexp *regexp.Regexp) []Asset {
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

var (
	imageRegexp      = regexp.MustCompile(`(?i)\[Image\(.*?image="(.*?)".*?\)]`)
	backgroundRegexp = regexp.MustCompile(`(?i)\[Background\(.*?image="(.*?)".*?\)]`)

	// characterRegexp matches the character directive in the story text.
	// 	Example: "[character(name="avg_npc_415_1#5$1",name2="char_011_talula_1#3",focus=2)]"
	//  group 1: the name (including # and $) of the 1st character, which is "avg_npc_415_1#5$1".
	//	group 2 (optional): the name (including # and $) of the 2nd character, which is "char_011_talula_1#3".
	characterRegexp = regexp.MustCompile(`(?i)(?U)\[Character\(.*name="(.*)".*(?:name2="(.*)".*)?\)\]`)

	// characterNameRegexp matches the name of a character.
	//
	// 	Example: "avg_103_angel_1#11$1"
	// 	group 1: the name of the character, which is "avg_103_angel_1".
	// 	group 2 (optional): the face number of the character (after the hash sign), which is "11".
	// 	group 3 (optional): the body number of the character (after the dollar sign), which is "1".
	characterNameRegexp = regexp.MustCompile(`^(.*?)(?:#(\d+))?(?:\$(\d+))?$`)
)

func findCharAssetsFromTxt(txt string) ([]Asset, error) {
	var assets []Asset
	matches := characterRegexp.FindAllStringSubmatch(txt, -1)
	for _, match := range matches {
		// match[1] is "name", required
		caName, err := normalizeCharAssetName(match[1])
		if err != nil {
			return nil, err
		}
		assets = append(assets, Asset{
			Name: caName,
			Kind: KindCharacter,
		})

		// match[2] is "name2", optional
		if match[2] != "" {
			caName, err := normalizeCharAssetName(match[2])
			if err != nil {
				return nil, err
			}
			assets = append(assets, Asset{
				Name: caName,
				Kind: KindCharacter,
			})
		}
	}
	return assets, nil
}

func normalizeCharAssetName(name string) (string, error) {
	matches := characterNameRegexp.FindStringSubmatch(name)
	if len(matches) == 0 {
		return "", fmt.Errorf("invalid character name: %v", name)
	}

	var err error
	faceNum, bodyNum := 1, 1 // default values are 1
	if matches[2] != "" {
		faceNum, err = strconv.Atoi(matches[2])
		if err != nil {
			return "", fmt.Errorf("invalid character face number: %v", matches[2])
		}
	}
	if matches[3] != "" {
		bodyNum, err = strconv.Atoi(matches[3])
		if err != nil {
			return "", fmt.Errorf("invalid character face number: %v", matches[3])
		}
	}

	name = matches[1]
	return fmt.Sprintf("%s#%d$%d", name, faceNum, bodyNum), nil
}
