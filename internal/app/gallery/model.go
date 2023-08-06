package gallery

import "github.com/flandiayingman/arkwaifu/internal/pkg/ark"

type Gallery struct {
	Server ark.Server `json:"-" gorm:"primaryKey;type:game_server"`

	ID          string `json:"id" gorm:"primaryKey;type:text COLLATE numeric;check:id=lower(id)"`
	Name        string `json:"name" gorm:""`
	Description string `json:"description" gorm:""`
	Arts        []Art  `json:"arts" gorm:"foreignKey:Server,GalleryID;references:Server,ID;constraint:OnDelete:CASCADE"`
}

func (g Gallery) TableName() string {
	return "gallery_galleries"
}

type Art struct {
	Server    ark.Server `json:"-" gorm:"primaryKey;type:game_server"`
	GalleryID string     `json:"-" gorm:""`
	SortID    int        `json:"-" gorm:""`

	ID          string `json:"id" gorm:"primaryKey;type:text COLLATE numeric;check:id=lower(id)"`
	Name        string `json:"name" gorm:""`
	Description string `json:"description" gorm:""`
	ArtID       string `json:"artID" gorm:"type:text COLLATE numeric;check:id=lower(id)"`
}

func (a Art) TableName() string {
	return "gallery_arts"
}
