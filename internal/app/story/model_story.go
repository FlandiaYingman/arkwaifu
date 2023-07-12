package story

import "github.com/flandiayingman/arkwaifu/internal/pkg/ark"

type Story struct {
	Server ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID     string     `json:"id" gorm:"primaryKey;collate:numeric;check:id=lower(id)"`

	Tag     Tag    `json:"tag" gorm:"type:story_tag"`
	TagText string `json:"tagText" gorm:""`
	Code    string `json:"code" gorm:"collate:numeric"`
	Name    string `json:"name" gorm:""`
	Info    string `json:"info" gorm:""`

	GroupID string  `json:"groupID" gorm:"collate:numeric"`
	SortID  *uint64 `json:"-" gorm:"unique;autoIncrement"`

	PictureArts   []PictureArt   `json:"pictureArts" gorm:"foreignKey:Server,StoryID;reference:(Server,ID);constraint:OnDelete:CASCADE"`
	CharacterArts []CharacterArt `json:"characterArts" gorm:"foreignKey:Server,StoryID;reference:(Server,ID);constraint:OnDelete:CASCADE"`
}

type Tag = string

const (
	TagBefore    Tag = "before"
	TagAfter     Tag = "after"
	TagInterlude Tag = "interlude"
)
