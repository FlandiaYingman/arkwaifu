package story

import "github.com/flandiayingman/arkwaifu/internal/pkg/ark"

type Group struct {
	Server  ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID      string     `json:"id" gorm:"primaryKey;check:id=lower(id)"`
	Name    string     `json:"name" gorm:""`
	Type    GroupType  `json:"type" gorm:"type:story_group_type"`
	Stories []Story    `json:"stories" gorm:"foreignKey:Server,GroupID;reference:Server,ID;constraint:OnDelete:CASCADE"`

	SortID *uint64 `json:"-" gorm:"unique;autoIncrement"`
}

type GroupType = string

const (
	GroupTypeMainStory  GroupType = "main-story"
	GroupTypeMajorEvent GroupType = "major-event"
	GroupTypeMinorEvent GroupType = "minor-event"
	GroupTypeOther      GroupType = "other"
)
