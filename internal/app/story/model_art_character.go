package story

import (
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/lib/pq"
)

type CharacterArt struct {
	Server   ark.Server `json:"server" gorm:"primaryKey;type:game_server"`
	ID       string     `json:"id" gorm:"primaryKey;check:id=lower(id)"`
	StoryID  string     `json:"storyID" gorm:"primaryKey"`
	Category string     `json:"category" gorm:""`

	Names pq.StringArray `json:"names" gorm:"type:text[];default:array[]::text[]"`

	SortID *uint64 `json:"-" gorm:"unique;autoIncrement"`
}

type AggregatedCharacterArt struct {
	Server   ark.Server `json:"server"`
	ID       string     `json:"id"`
	Category string     `json:"category"`

	Names pq.StringArray `json:"names" gorm:"type:text[];default:array[]::text[]"`
}
