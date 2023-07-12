package story

import "github.com/flandiayingman/arkwaifu/internal/pkg/ark"

type PictureArt struct {
	Server   ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID       string     `json:"id" gorm:"primaryKey;check:id=lower(id)"`
	StoryID  string     `json:"storyID" gorm:"primaryKey"`
	Category string     `json:"category" gorm:""`

	Title    string `json:"title" gorm:""`
	Subtitle string `json:"subtitle" gorm:""`

	SortID *uint64 `json:"-" gorm:"unique;autoIncrement"`
}
