package story

import (
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/lib/pq"
)

type Story struct {
	Server ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID     string     `json:"id" gorm:"primaryKey;collate:numeric;check:id=lower(id)"`

	Tag     Tag    `json:"tag" gorm:"type:story_tag"`
	TagText string `json:"tagText" gorm:""`
	Code    string `json:"code" gorm:"collate:numeric"`
	Name    string `json:"name" gorm:""`
	Info    string `json:"info" gorm:""`

	GroupID string `json:"groupID" gorm:"collate:numeric"`

	PictureArts   []PictureArt   `json:"pictureArts" gorm:"foreignKey:Server,StoryID;reference:(Server,ID)"`
	CharacterArts []CharacterArt `json:"characterArts" gorm:"foreignKey:Server,StoryID;reference:(Server,ID)"`
}

type Tag = string

const (
	TagBefore    Tag = "before"
	TagAfter     Tag = "after"
	TagInterlude Tag = "interlude"
)

type Group struct {
	Server  ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID      string     `json:"id" gorm:"primaryKey;collate:numeric;check:id=lower(id)"`
	Name    string     `json:"name" gorm:""`
	Type    GroupType  `json:"type" gorm:"type:story_group_type"`
	Stories []Story    `json:"stories" gorm:"foreignKey:Server,GroupID;reference:Server,ID"`
}

type GroupType = string

const (
	GroupTypeMainStory  GroupType = "main-story"
	GroupTypeMajorEvent GroupType = "major-event"
	GroupTypeMinorEvent GroupType = "minor-event"
	GroupTypeOther      GroupType = "other"
)

type Tree = []Group

type PictureArt struct {
	Server   ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID       string     `json:"id" gorm:"primaryKey;check:id=lower(id)"`
	StoryID  string     `json:"storyID" gorm:"primaryKey"`
	Category string     `json:"category" gorm:""`

	Title    string `json:"title" gorm:""`
	Subtitle string `json:"subtitle" gorm:""`
}

type CharacterArt struct {
	Server   ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID       string     `json:"id" gorm:"primaryKey;check:id=lower(id)"`
	StoryID  string     `json:"storyID" gorm:"primaryKey"`
	Category string     `json:"category" gorm:""`

	Names pq.StringArray `json:"names" gorm:"type:text[]"`
}
