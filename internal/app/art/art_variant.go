package art

import (
	"fmt"
)

type Variant struct {
	ArtID     string `gorm:"primaryKey;type:text COLLATE numeric;check:art_id=lower(art_id)" json:"artID"`
	Variation string `gorm:"primaryKey" json:"variation,omitempty"`

	ContentPresent bool `gorm:"" json:"contentPresent"`
	ContentWidth   *int `gorm:"" json:"contentWidth,omitempty"`
	ContentHeight  *int `gorm:"" json:"contentHeight,omitempty"`
}

const (
	VariationOrigin     string = "origin"
	VariationThumbnail  string = "thumbnail"
	VariationRealEsrgan string = "real-esrgan"
)

func NewVariant(id string, variation string) *Variant {
	return &Variant{
		ArtID:          id,
		Variation:      variation,
		ContentPresent: false,
		ContentWidth:   nil,
		ContentHeight:  nil,
	}
}

func (v *Variant) String() string {
	return fmt.Sprintf("%s/%s", v.Variation, v.ArtID)
}
