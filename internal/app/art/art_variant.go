package art

import (
	"fmt"
)

type Variant struct {
	ArtID     string    `gorm:"primaryKey;type:text COLLATE numeric;check:art_id=lower(art_id)" json:"artID"`
	Variation Variation `gorm:"primaryKey;type:art_variation" json:"variation,omitempty"`

	ContentPresent bool `gorm:"" json:"contentPresent"`
	ContentWidth   *int `gorm:"" json:"contentWidth,omitempty"`
	ContentHeight  *int `gorm:"" json:"contentHeight,omitempty"`
}

func NewVariant(id string, variation Variation) *Variant {
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

type Variation string

const (
	VariationOrigin                Variation = "origin"
	VariationRealEsrganX4Plus      Variation = "real-esrgan(realesrgan-x4plus)"
	VariationRealEsrganX4PlusAnime Variation = "real-esrgan(realesrgan-x4plus-anime)"
	VariationThumbnail             Variation = "thumbnail"
)

func ParseVariation(str string) (Variation, error) {
	variation := Variation(str)
	switch variation {
	case
		VariationOrigin,
		VariationRealEsrganX4Plus,
		VariationRealEsrganX4PlusAnime,
		VariationThumbnail:
		return variation, nil
	default:
		return "", fmt.Errorf("string %q is not a variation", str)
	}
}

func (v *Variation) UnmarshalJSON(bytes []byte) error {
	variation, err := ParseVariation(string(bytes))
	if err != nil {
		return err
	}
	*v = variation
	return nil
}
