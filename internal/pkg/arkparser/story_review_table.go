package arkparser

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wk8/go-ordered-map/v2"
	"os"
	"path"
	"path/filepath"
)

type JsonStoryGroup struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    string       `json:"actType"`
	Stories []*JsonStory `json:"infoUnlockDatas"`
}
type JsonStory struct {
	ID      string       `json:"storyId"`
	GroupID string       `json:"storyGroup"`
	Code    string       `json:"storyCode"`
	Name    string       `json:"storyName"`
	Info    string       `json:"storyInfo"`
	Text    string       `json:"storyTxt"`
	Tag     JsonStoryTag `json:"avgTag"`
}
type JsonStoryTag struct {
	Type string
	Text string
}

func (t *JsonStoryTag) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return errors.WithStack(err)
	}
	switch s {
	case TagBeforeCN, TagBeforeEN, TagBeforeJP, TagBeforeKR:
		*t = JsonStoryTag{
			Type: TagBefore,
			Text: s,
		}
	case TagAfterCN, TagAfterEN, TagAfterJP, TagAfterKR:
		*t = JsonStoryTag{
			Type: TagAfter,
			Text: s,
		}
	case TagInterludeCN, TagInterludeEN, TagInterludeJP, TagInterludeKR:
		*t = JsonStoryTag{
			Type: TagInterlude,
			Text: s,
		}
	default:
		return errors.Errorf("unknown story tag: %s", s)
	}
	return nil
}

const (
	TypeMain         string = "MAIN_STORY"
	TypeActivity     string = "ACTIVITY_STORY"
	TypeMiniActivity string = "MINI_STORY"
	TypeNone         string = "NONE"

	TagBefore    string = "BEFORE"
	TagAfter     string = "AFTER"
	TagInterlude string = "INTERLUDE"

	TagBeforeCN    string = "行动前"
	TagAfterCN     string = "行动后"
	TagInterludeCN string = "幕间"
	TagBeforeEN    string = "Before Operation"
	TagAfterEN     string = "After Operation"
	TagInterludeEN string = "Interlude"
	TagBeforeJP    string = "戦闘前"
	TagAfterJP     string = "戦闘後"
	TagInterludeJP string = "幕間"
	TagBeforeKR    string = "작전 전"
	TagAfterKR     string = "작전 후"
	TagInterludeKR string = "브릿지"
)

func (p *Parser) ParseStoryReviewTable() ([]*JsonStoryGroup, error) {
	jsonPath := filepath.Join(p.Root, p.Prefix, "gamedata/excel/story_review_table.json")
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// key: story group id; value: story group
	dynamicJsonObject := orderedmap.New[string, JsonStoryGroup]()
	err = json.Unmarshal(jsonData, &dynamicJsonObject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*JsonStoryGroup, 0, dynamicJsonObject.Len())
	for pair := dynamicJsonObject.Oldest(); pair != nil; pair = pair.Next() {
		result = append(result, &pair.Value)
	}

	return result, nil
}

func (p *Parser) GetInfo(infoPath string) (string, error) {
	infoPath = fmt.Sprintf("[uc]%s.txt", infoPath)
	infoPath = path.Join(p.Root, p.Prefix, "gamedata/story", infoPath)

	bytes, err := os.ReadFile(infoPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(bytes), nil
}
func (p *Parser) GetText(textPath string) (string, error) {
	textPath = fmt.Sprintf("%s.txt", textPath)
	textPath = path.Join(p.Root, p.Prefix, "gamedata/story", textPath)

	bytes, err := os.ReadFile(textPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(bytes), nil
}
