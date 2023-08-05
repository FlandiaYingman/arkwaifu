package arkjson

import (
	"github.com/Jeffail/gabs/v2"
	"path/filepath"
)

const (
	StoryReviewMetaTablePath = "assets/torappu/dynamicassets/gamedata/excel/story_review_meta_table.json"
	ReplicateTable           = "assets/torappu/dynamicassets/gamedata/excel/replicate_table.json"
	RetroTable               = "assets/torappu/dynamicassets/gamedata/excel/retro_table.json"
	RoguelikeTopicTable      = "assets/torappu/dynamicassets/gamedata/excel/roguelike_topic_table.json"
)

func Get(root string, path string) (*gabs.Container, error) {
	json, err := gabs.ParseJSONFile(filepath.Join(root, path))
	if err != nil {
		return nil, err
	}
	return json, nil
}
