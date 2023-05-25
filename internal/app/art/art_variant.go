package art

import (
	"fmt"
)

type Variant struct {
	ArtID     string `gorm:"primaryKey;type:text COLLATE numeric" json:"artID"`
	Variation string `gorm:"primaryKey" json:"variation,omitempty"`

	ContentPresent bool   `gorm:"" json:"contentPresent,omitempty"`
	ContentPath    string `gorm:"" json:"contentPath,omitempty"`
	ContentWidth   *int   `gorm:"" json:"contentWidth,omitempty"`
	ContentHeight  *int   `gorm:"" json:"contentHeight,omitempty"`
}

const (
	VariationOrigin     string = "origin"
	VariationThumbnail  string = "thumbnail"
	VariationRealEsrgan string = "real-esrgan"
)

func NewVariant(id string, variation string) *Variant {
	return &Variant{
		ArtID:     id,
		Variation: variation,
	}
}

func (v *Variant) String() string {
	return fmt.Sprintf("%s/%s", v.Variation, v.ArtID)
}

func (v *Variant) ToStatic() *VariantContent {
	return &VariantContent{
		ArtID:     v.ArtID,
		Variation: v.Variation,
		Content:   nil,
	}
}
