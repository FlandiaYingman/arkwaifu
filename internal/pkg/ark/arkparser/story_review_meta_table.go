package arkparser

import (
	"encoding/json"
	"github.com/wk8/go-ordered-map/v2"
	"os"
	"path/filepath"
)

type JsonResPic struct {
	Id             string `json:"id"`
	Desc           string `json:"desc"`
	AssetPath      string `json:"assetPath"`
	PicDescription string `json:"picDescription"`
}

type JsonComponent struct {
	Pic struct {
		Pics orderedmap.OrderedMap[string, JsonPic] `json:"pics"`
	} `json:"pic"`
}

type JsonPic struct {
	PicId     string `json:"picId"`
	PicSortId int    `json:"picSortId"`
}

type JsonStoryReviewMetaTable struct {
	ActArchiveResData struct {
		Pics orderedmap.OrderedMap[string, JsonResPic] `json:"pics"`
	} `json:"actArchiveResData"`
	ActArchiveData struct {
		Components orderedmap.OrderedMap[string, JsonComponent] `json:"components"`
	} `json:"actArchiveData"`
}

func (p *Parser) ParseStoryReviewMetaTable() (*JsonStoryReviewMetaTable, error) {
	jsonPath := filepath.Join(p.Root, p.Prefix, "gamedata/excel/story_review_meta_table.json")
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	var obj JsonStoryReviewMetaTable
	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
